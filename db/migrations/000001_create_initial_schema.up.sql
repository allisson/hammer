-- topics table

CREATE TABLE IF NOT EXISTS topics(
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS created_at_idx ON topics USING BRIN (created_at);

-- subscriptions table

CREATE TABLE IF NOT EXISTS subscriptions(
    id VARCHAR PRIMARY KEY,
    topic_id VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    url VARCHAR NOT NULL,
    secret_token VARCHAR NOT NULL,
    max_delivery_attempts INT NOT NULL,
    delivery_attempt_delay INT NOT NULL,
    delivery_attempt_timeout INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (topic_id) REFERENCES topics (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS topic_id_idx ON subscriptions (topic_id);
CREATE INDEX IF NOT EXISTS created_at_idx ON subscriptions USING BRIN (created_at);

-- messages table

CREATE TABLE IF NOT EXISTS messages(
    id VARCHAR PRIMARY KEY,
    topic_id VARCHAR NOT NULL,
    content_type VARCHAR NOT NULL,
    data TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (topic_id) REFERENCES topics (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS topic_id_idx ON messages (topic_id);
CREATE INDEX IF NOT EXISTS created_at_idx ON messages USING BRIN (created_at);

-- deliveries table

CREATE TABLE IF NOT EXISTS deliveries(
    id VARCHAR PRIMARY KEY,
    topic_id VARCHAR NOT NULL,
    subscription_id VARCHAR NOT NULL,
    message_id VARCHAR NOT NULL,
    content_type VARCHAR NOT NULL,
    data TEXT NOT NULL,
    url VARCHAR NOT NULL,
    secret_token VARCHAR NOT NULL,
    max_delivery_attempts INT NOT NULL,
    delivery_attempt_delay INT NOT NULL,
    delivery_attempt_timeout INT NOT NULL,
    scheduled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    delivery_attempts INT NOT NULL,
    status VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (topic_id) REFERENCES topics (id) ON DELETE CASCADE,
    FOREIGN KEY (subscription_id) REFERENCES subscriptions (id) ON DELETE CASCADE,
    FOREIGN KEY (message_id) REFERENCES messages (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS topic_id_idx ON deliveries (topic_id);
CREATE INDEX IF NOT EXISTS subscription_id_idx ON deliveries (subscription_id);
CREATE INDEX IF NOT EXISTS message_id_idx ON deliveries (message_id);
CREATE INDEX IF NOT EXISTS scheduled_at_idx ON deliveries USING BRIN (scheduled_at);
CREATE INDEX IF NOT EXISTS status_idx ON deliveries (status);
CREATE INDEX IF NOT EXISTS created_at_idx ON deliveries USING BRIN (created_at);

-- delivery_attempts table

CREATE TABLE IF NOT EXISTS delivery_attempts(
    id VARCHAR PRIMARY KEY,
    delivery_id VARCHAR NOT NULL,
    request TEXT NOT NULL,
    response TEXT NOT NULL,
    response_status_code INT NOT NULL,
    execution_duration INT NOT NULL,
    success BOOLEAN NOT NULL,
    error TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (delivery_id) REFERENCES deliveries (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS delivery_id_idx ON delivery_attempts (delivery_id);
CREATE INDEX IF NOT EXISTS created_at_idx ON delivery_attempts USING BRIN (created_at);
