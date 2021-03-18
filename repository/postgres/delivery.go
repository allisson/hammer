package repository

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/allisson/hammer"
	"github.com/jmoiron/sqlx"
)

type dispatchResponse struct {
	Request            string
	Response           string
	ResponseStatusCode int
	ExecutionDuration  int
	Success            bool
	Error              string
}

func dispatchToURL(delivery *hammer.Delivery) dispatchResponse {
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
	httpClient := &http.Client{Timeout: time.Duration(delivery.DeliveryAttemptTimeout) * time.Second}
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

// Delivery is a implementation of hammer.DeliveryRepository
type Delivery struct {
	db *sqlx.DB
}

// Find returns hammer.Delivery by id
func (d Delivery) Find(ctx context.Context, id string) (*hammer.Delivery, error) {
	delivery := &hammer.Delivery{}
	findOptions := hammer.FindOptions{
		FindFilters: []hammer.FindFilter{
			{
				FieldName: "id",
				Operator:  "=",
				Value:     id,
			},
		},
	}
	query, args := findQuery("deliveries", findOptions)
	err := d.db.GetContext(ctx, delivery, query, args...)
	return delivery, err
}

// FindAll returns []hammer.Delivery by limit and offset
func (d Delivery) FindAll(ctx context.Context, findOptions hammer.FindOptions) ([]*hammer.Delivery, error) {
	deliveries := []*hammer.Delivery{}
	query, args := findQuery("deliveries", findOptions)
	err := d.db.SelectContext(ctx, &deliveries, query, args...)
	return deliveries, err
}

// Store a hammer.Delivery on database (create or update)
func (d Delivery) Store(ctx context.Context, delivery *hammer.Delivery) error {
	_, err := d.Find(ctx, delivery.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			query, args := insertQuery("deliveries", delivery)
			_, err := d.db.ExecContext(ctx, query, args...)
			return err
		}
		return err
	}
	query, args := updateQuery("deliveries", delivery.ID, delivery)
	_, err = d.db.ExecContext(ctx, query, args...)
	return err
}

// Dispatch fetchs a delivery and send to url destination.
func (d Delivery) Dispatch(ctx context.Context) (*hammer.DeliveryAttempt, error) {
	// Generate delivery attempt id
	id, err := hammer.GenerateULID()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			*
		FROM
			deliveries
		WHERE
			deliveries.status = $1 AND deliveries.scheduled_at <= $2
		ORDER BY
			deliveries.created_at ASC
		FOR UPDATE SKIP LOCKED
		LIMIT
			1
	`

	// Starts a new transaction
	tx, err := d.db.Beginx()
	if err != nil {
		return nil, err
	}

	// Get delivery
	delivery := hammer.Delivery{}
	err = tx.GetContext(ctx, &delivery, query, hammer.DeliveryStatusPending, time.Now().UTC())
	if err != nil {
		// Skip if no result
		if err == sql.ErrNoRows {
			rollback(tx)
			return nil, nil
		}
		rollback(tx)
		return nil, err
	}

	// Dispatch webhook
	dr := dispatchToURL(&delivery)

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
	query, args := insertQuery("delivery_attempts", deliveryAttempt)
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		rollback(tx)
		return nil, err
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
	query, args = updateQuery("deliveries", delivery.ID, delivery)
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		rollback(tx)
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &deliveryAttempt, nil
}

// NewDelivery returns a new Delivery with db connection
func NewDelivery(db *sqlx.DB) *Delivery {
	return &Delivery{db: db}
}
