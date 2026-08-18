CREATE TABLE events(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    event_type VARCHAR(20) NOT NULL,
    user_id BINARY(16) NOT NULL,
    central_session_id BINARY(16),
    application_id BINARY(16),
    payload JSON NOT NULL,
    status VARCHAR(20),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    published_at DATETIME NOT NULL
);

CREATE TABLE event_deliveries(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    event_id BINARY(16) NOT NULL,
    application_id BINARY(16) NOT NULL,
    status VARCHAR(20) NOT NULL,
    attempt_count INTEGER DEFAULT 0 NOT NULL,
    last_attempt_at DATETIME,
    next_retry_at DATETIME,
    processed_at DATETIME,
    last_error TEXT,
    CONSTRAINT fk_event_deliv
        FOREIGN KEY (event_id)
        REFERENCES events(id)
        ON DELETE CASCADE
)