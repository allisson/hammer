package repository

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"text/template"
)

var queriesCache sync.Map

const (
	// Generic Statements
	sqlFind = `
		SELECT *
		FROM {{.table}}
		WHERE id = @id
	`
	sqlFindAll = `
		SELECT *
		FROM {{.table}}
		{{if .orderBy}}ORDER BY @orderBy{{end}}
		{{if .limit}}LIMIT @limit{{end}}
		{{if .offset}}OFFSET @offset{{end}}
	`
	// DeliveryAttempt Statements
	sqlDeliveryAttemptCreate = `
		INSERT INTO delivery_attempts (
			"id",
			"delivery_id",
			"request",
			"response",
			"response_status_code",
			"execution_duration",
			"success",
			"error",
			"created_at"
		)
		VALUES (
			:id,
			:delivery_id,
			:request,
			:response,
			:response_status_code,
			:execution_duration,
			:success,
			:error,
			:created_at
		)
	`
	sqlDeliveryAttemptUpdate = `
		UPDATE delivery_attempts
		SET delivery_id = :delivery_id,
			request = :request,
			response = :response,
			response_status_code = :response_status_code,
			execution_duration = :execution_duration,
			success = :success,
			error = :error,
			created_at = :created_at
		WHERE id = :id
	`
	// Delivery Statements
	sqlDeliveryFindToDispatch = `
		SELECT id
		FROM deliveries
		WHERE status = $1 AND scheduled_at < $2
		ORDER BY id ASC
		LIMIT $3
		OFFSET $4
	`
	sqlDeliveryCreate = `
		INSERT INTO deliveries (
			"id",
			"topic_id",
			"subscription_id",
			"message_id",
			"data",
			"url",
			"secret_token",
			"max_delivery_attempts",
			"delivery_attempt_delay",
			"delivery_attempt_timeout",
			"scheduled_at",
			"delivery_attempts",
			"status",
			"created_at",
			"updated_at"
		)
		VALUES (
			:id,
			:topic_id,
			:subscription_id,
			:message_id,
			:data,
			:url,
			:secret_token,
			:max_delivery_attempts,
			:delivery_attempt_delay,
			:delivery_attempt_timeout,
			:scheduled_at,
			:delivery_attempts,
			:status,
			:created_at,
			:updated_at
		)
	`
	sqlDeliveryUpdate = `
		UPDATE deliveries
		SET topic_id = :topic_id,
			subscription_id = :subscription_id,
			message_id = :message_id,
			data = :data,
			url = :url,
			secret_token = :secret_token,
			max_delivery_attempts = :max_delivery_attempts,
			delivery_attempt_delay = :delivery_attempt_delay,
			delivery_attempt_timeout = :delivery_attempt_timeout,
			scheduled_at = :scheduled_at,
			delivery_attempts = :delivery_attempts,
			status = :status,
			created_at = :created_at,
			updated_at = :updated_at
		WHERE id = :id
	`
	// Message Statements
	sqlMessageFindAll = `
		SELECT *
		FROM messages
		{{if .topicID}}WHERE topic_id = @topic_id{{end}}
		{{if .orderBy}}ORDER BY @orderBy{{end}}
		{{if .limit}}LIMIT @limit{{end}}
		{{if .offset}}OFFSET @offset{{end}}
	`
	sqlMessageCreate = `
		INSERT INTO messages (
			"id",
			"topic_id",
			"data",
			"created_at"
		)
		VALUES (
			:id,
			:topic_id,
			:data,
			:created_at
		)
	`
	sqlMessageUpdate = `
		UPDATE messages
		SET topic_id = :topic_id,
			data = :data,
			created_at = :created_at
		WHERE id = :id
	`
	// Subscription Statements
	sqlSubscriptionFind = `
		SELECT *
		FROM subscriptions
		WHERE id = $1
	`
	sqlSubscriptionCreate = `
		INSERT INTO subscriptions (
			"id",
			"topic_id",
			"name",
			"url",
			"secret_token",
			"max_delivery_attempts",
			"delivery_attempt_delay",
			"delivery_attempt_timeout",
			"created_at",
			"updated_at"
		)
		VALUES (
			:id,
			:topic_id,
			:name,
			:url,
			:secret_token,
			:max_delivery_attempts,
			:delivery_attempt_delay,
			:delivery_attempt_timeout,
			:created_at,
			:updated_at
		)
	`
	sqlSubscriptionUpdate = `
		UPDATE subscriptions
		SET topic_id = :topic_id,
			name = :name,
			url = :url,
			secret_token = :secret_token,
			max_delivery_attempts = :max_delivery_attempts,
			delivery_attempt_delay = :delivery_attempt_delay,
			delivery_attempt_timeout = :delivery_attempt_timeout,
			created_at = :created_at,
			updated_at = :updated_at
		WHERE id = :id
	`
	// Topic Statements
	sqlTopicCreate = `
		INSERT INTO topics (
			"id",
			"name",
			"created_at",
			"updated_at"
		)
		VALUES (
			:id,
			:name,
			:created_at,
			:updated_at
		)
	`
	sqlTopicUpdate = `
		UPDATE topics
		SET name = :name,
			created_at = :created_at,
			updated_at = :updated_at
		WHERE id = :id
	`
)

func buildQuery(text string, data map[string]interface{}) (string, []interface{}, error) {
	var t *template.Template
	v, ok := queriesCache.Load(text)
	if !ok {
		var err error
		t, err = template.New("query").Parse(text)
		if err != nil {
			return "", nil, fmt.Errorf("could not parse sql query template: %w", err)
		}

		queriesCache.Store(text, t)
	} else {
		t = v.(*template.Template)
	}

	var wr bytes.Buffer
	if err := t.Execute(&wr, data); err != nil {
		return "", nil, fmt.Errorf("could not apply sql query data: %w", err)
	}

	query := wr.String()
	args := []interface{}{}
	for key, val := range data {
		if !strings.Contains(query, "@"+key) {
			continue
		}

		args = append(args, val)
		query = strings.Replace(query, "@"+key, fmt.Sprintf("$%d", len(args)), -1)
	}
	return query, args, nil
}
