package repository

const (
	// DeliveryAttempt Statements
	sqlDeliveryAttemptFind = `
		SELECT *
		FROM delivery_attempts
		WHERE id = $1
	`
	sqlDeliveryAttemptFindAll = `
		SELECT *
		FROM delivery_attempts
		ORDER BY id DESC
		LIMIT $1
		OFFSET $2
	`
	sqlDeliveryAttemptCreate = `
		INSERT INTO delivery_attempts (
			"id",
			"delivery_id",
			"url",
			"request_headers",
			"request_body",
			"response_headers",
			"response_body",
			"response_status_code",
			"execution_duration",
			"success",
			"created_at"
		)
		VALUES (
			:id,
			:delivery_id,
			:url,
			:request_headers,
			:request_body,
			:response_headers,
			:response_body,
			:response_status_code,
			:execution_duration,
			:success,
			:created_at
		)
	`
	sqlDeliveryAttemptUpdate = `
		UPDATE delivery_attempts
		SET delivery_id = :delivery_id,
			url = :url,
			request_headers = :request_headers,
			request_body = :request_body,
			response_headers = :response_headers,
			response_body = :response_body,
			response_status_code = :response_status_code,
			execution_duration = :execution_duration,
			success = :success,
			created_at = :created_at
		WHERE id = :id
	`
	// Delivery Statements
	sqlDeliveryFind = `
		SELECT *
		FROM deliveries
		WHERE id = $1
	`
	sqlDeliveryFindAll = `
		SELECT *
		FROM deliveries
		ORDER BY id DESC
		LIMIT $1
		OFFSET $2
	`
	sqlDeliveryCreate = `
		INSERT INTO deliveries (
			"id",
			"topic_id",
			"subscription_id",
			"message_id",
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
			scheduled_at = :scheduled_at,
			delivery_attempts = :delivery_attempts,
			status = :status,
			created_at = :created_at,
			updated_at = :updated_at
		WHERE id = :id
	`
	// Message Statements
	sqlMessageFind = `
		SELECT *
		FROM messages
		WHERE id = $1
	`
	sqlMessageFindAll = `
		SELECT *
		FROM messages
		ORDER BY id DESC
		LIMIT $1
		OFFSET $2
	`
	sqlMessageFindByTopic = `
		SELECT *
		FROM messages
		WHERE topic_id = $1
		ORDER BY id DESC
		LIMIT $2
		OFFSET $3
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
	sqlSubscriptionFindAll = `
		SELECT *
		FROM subscriptions
		ORDER BY id ASC
		LIMIT $1
		OFFSET $2
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
	sqlTopicFind = `
		SELECT *
		FROM topics
		WHERE id = $1
	`
	sqlTopicFindAll = `
		SELECT *
		FROM topics
		ORDER BY id ASC
		LIMIT $1
		OFFSET $2
	`
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
