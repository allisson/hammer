package repository

import (
	"github.com/allisson/hammer"
	"github.com/huandu/go-sqlbuilder"
)

const (
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

func buildSQLQuery(tableName string, findOptions hammer.FindOptions) (sql string, args []interface{}) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").From(tableName)

	// Pagination
	if findOptions.FindPagination != nil {
		sb.Limit(int(findOptions.FindPagination.Limit)).Offset(int(findOptions.FindPagination.Offset))
	}

	// Filters
	for _, findFilter := range findOptions.FindFilters {
		switch findFilter.Operator {
		case "=":
			sb.Where(sb.Equal(findFilter.FieldName, findFilter.Value))
		case "gt":
			sb.Where(sb.GreaterThan(findFilter.FieldName, findFilter.Value))
		case "gte":
			sb.Where(sb.GreaterEqualThan(findFilter.FieldName, findFilter.Value))
		case "lt":
			sb.Where(sb.LessThan(findFilter.FieldName, findFilter.Value))
		case "lte":
			sb.Where(sb.LessEqualThan(findFilter.FieldName, findFilter.Value))
		}
	}

	// Order by
	if findOptions.FindOrderBy != nil {
		sb.OrderBy(findOptions.FindOrderBy.FieldName)
	}

	// Return result
	return sb.Build()
}
