package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/allisson/hammer"
)

type dispatchResponse struct {
	Request            string
	Response           string
	ResponseStatusCode int
	ExecutionDuration  int
	Success            bool
	Error              string
}

func makeRequest(delivery *hammer.Delivery, httpClient *http.Client) dispatchResponse {
	dr := dispatchResponse{}

	// Create payload
	cloudEvent := hammer.CloudEventPayload{
		SpecVersion:     "1.0",
		Type:            "hammer.message.created",
		Source:          fmt.Sprintf("/v1/messages/%s", delivery.MessageID),
		ID:              delivery.ID,
		Time:            delivery.CreatedAt,
		SecretToken:     delivery.SecretToken,
		MessageID:       delivery.MessageID,
		SubscriptionID:  delivery.SubscriptionID,
		TopicID:         delivery.TopicID,
		DataContentType: delivery.ContentType,
		DataBase64:      delivery.Data,
	}

	// Convert to json
	requestBody, err := json.Marshal(&cloudEvent)
	if err != nil {
		dr.Error = err.Error()
		return dr
	}

	// Prepare request
	request, err := http.NewRequest("POST", delivery.URL, bytes.NewBuffer(requestBody))
	if err != nil {
		dr.Error = err.Error()
		return dr
	}
	request.Header.Set("Content-Type", "application/json")
	requestDump, err := httputil.DumpRequest(request, true)
	if err != nil {
		dr.Error = err.Error()
		return dr
	}
	dr.Request = string(requestDump)

	// Make request
	start := time.Now()
	response, err := httpClient.Do(request)
	if err != nil {
		dr.Error = err.Error()
		return dr
	}
	latency := time.Since(start)
	responseDump, err := httputil.DumpResponse(response, true)
	if err != nil {
		dr.Error = err.Error()
		return dr
	}
	dr.Response = string(responseDump)

	// Update dispatch response
	dr.ResponseStatusCode = response.StatusCode
	dr.ExecutionDuration = int(latency.Milliseconds())
	switch dr.ResponseStatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent:
		dr.Success = true
	}

	return dr
}

// Delivery is a implementation of hammer.DeliveryService
type Delivery struct {
	deliveryRepo        hammer.DeliveryRepository
	deliveryAttemptRepo hammer.DeliveryAttemptRepository
	txFactoryRepo       hammer.TxFactoryRepository
}

// Find returns hammer.Delivery by id
func (d *Delivery) Find(id string) (hammer.Delivery, error) {
	return d.deliveryRepo.Find(id)
}

// FindAll returns []hammer.Delivery by findOptions
func (d *Delivery) FindAll(findOptions hammer.FindOptions) ([]hammer.Delivery, error) {
	return d.deliveryRepo.FindAll(findOptions)
}

// FindToDispatch returns []string ready to dispatch by limit and offset
func (d *Delivery) FindToDispatch(limit, offset int) ([]string, error) {
	return d.deliveryRepo.FindToDispatch(limit, offset)
}

// Dispatch message to destination
func (d *Delivery) Dispatch(delivery *hammer.Delivery, httpClient *http.Client) error {
	// Generate delivery attempt id
	id, err := generateULID()
	if err != nil {
		return err
	}

	dr := makeRequest(delivery, httpClient)

	// Start tx
	tx, err := d.txFactoryRepo.New()
	if err != nil {
		return err
	}

	// Create delivery attempt
	deliveryAttempt := hammer.DeliveryAttempt{
		ID:                 id,
		DeliveryID:         delivery.ID,
		Request:            dr.Request,
		Response:           dr.Response,
		ResponseStatusCode: dr.ResponseStatusCode,
		ExecutionDuration:  dr.ExecutionDuration,
		Success:            dr.Success,
		Error:              dr.Error,
		CreatedAt:          time.Now().UTC(),
	}
	err = d.deliveryAttemptRepo.Store(tx, &deliveryAttempt)
	if err != nil {
		rollback(tx, "delivery-dispatch-delivery-attempt-store")
		return err
	}

	// Update delivery
	delivery.DeliveryAttempts++
	delivery.UpdatedAt = time.Now().UTC()
	if deliveryAttempt.Success {
		delivery.Status = hammer.DeliveryStatusCompleted
	} else {
		if delivery.DeliveryAttempts >= delivery.MaxDeliveryAttempts {
			delivery.Status = hammer.DeliveryStatusFailed
		} else {
			delivery.ScheduledAt = time.Now().UTC().Add(time.Duration(delivery.DeliveryAttemptDelay) * time.Second)
		}
	}
	err = d.deliveryRepo.Store(tx, delivery)
	if err != nil {
		rollback(tx, "delivery-dispatch-delivery-store")
		return err
	}

	// Commit tx
	err = tx.Commit()
	if err != nil {
		rollback(tx, "delivery-dispatch-commit")
		return err
	}

	return nil
}

// NewDelivery returns a new Delivery with DeliveryRepo
func NewDelivery(deliveryRepo hammer.DeliveryRepository, deliveryAttemptRepo hammer.DeliveryAttemptRepository, txFactoryRepo hammer.TxFactoryRepository) Delivery {
	return Delivery{
		deliveryRepo:        deliveryRepo,
		deliveryAttemptRepo: deliveryAttemptRepo,
		txFactoryRepo:       txFactoryRepo,
	}
}
