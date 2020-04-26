-- topics table

CREATE TABLE IF NOT EXISTS topics(
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    active BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- subscriptions table

CREATE TABLE IF NOT EXISTS subscriptions(
    id VARCHAR PRIMARY KEY,
    topic_id VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    url VARCHAR NOT NULL,
    secret_token VARCHAR NOT NULL,
    active BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (topic_id) REFERENCES topics (id)
);

CREATE INDEX IF NOT EXISTS topic_id_idx ON subscriptions (topic_id);
CREATE INDEX IF NOT EXISTS created_at_idx ON subscriptions (created_at);

-- messages table

CREATE TABLE IF NOT EXISTS messages(
    id VARCHAR PRIMARY KEY,
    topic_id VARCHAR NOT NULL,
    data TEXT NOT NULL,
    created_deliveries BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (topic_id) REFERENCES topics (id)
);

CREATE INDEX IF NOT EXISTS topic_id_idx ON messages (topic_id);
CREATE INDEX IF NOT EXISTS created_deliveries_idx ON messages (created_deliveries);

-- deliveries table

CREATE TABLE IF NOT EXISTS deliveries(
    id VARCHAR PRIMARY KEY,
    topic_id VARCHAR NOT NULL,
    subscription_id VARCHAR NOT NULL,
    message_id VARCHAR NOT NULL,
    max_delivery_attempts INT NOT NULL,
    delivery_attempt_delay INT NOT NULL,
    delivery_attempt_timeout INT NOT NULL,
    delivery_attempts INT NOT NULL,
    last_delivery_attempt TIMESTAMPTZ,
    status VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (topic_id) REFERENCES topics (id),
    FOREIGN KEY (subscription_id) REFERENCES subscriptions (id),
    FOREIGN KEY (message_id) REFERENCES messages (id)
);

CREATE INDEX IF NOT EXISTS topic_id_idx ON deliveries (topic_id);
CREATE INDEX IF NOT EXISTS subscription_id_idx ON deliveries (subscription_id);
CREATE INDEX IF NOT EXISTS message_id_idx ON deliveries (message_id);

-- delivery_attemps table

CREATE TABLE IF NOT EXISTS delivery_attemps(
    id VARCHAR PRIMARY KEY,
    delivery_id VARCHAR NOT NULL,
    url VARCHAR NOT NULL,
    request_headers TEXT NOT NULL,
    request_body TEXT NOT NULL,
    response_headers TEXT NOT NULL,
    response_body TEXT NOT NULL,
    response_status_code INT NOT NULL,
    execution_duration INT NOT NULL,
    success BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (delivery_id) REFERENCES deliveries (id)
);

CREATE INDEX IF NOT EXISTS delivery_id_idx ON delivery_attemps (delivery_id);
