CREATE TABLE code_verifiers(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    state VARCHAR(32) NOT NULL UNIQUE,
    code_verifier VARCHAR(64) NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE profile_cache(
    external_user_id BINARY(16) PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    email VARCHAR(50) NOT NULL,
    groups_list JSON,
    synced_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP   
);

CREATE TABLE local_sessions(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    session_token_hash VARCHAR(64) NOT NULL,
    external_user_id BINARY(16) NOT NULL,
    central_session_id BINARY(16) NOT NULL,
    application_id BINARY(16),
    status ENUM('Active', 'Expired', 'Revoked') NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at DATETIME NOT NULL,
    last_activity_at DATETIME,
    revoked_at DATETIME,
    revoked_reason VARCHAR(100),
    CONSTRAINT fk_session_user
        FOREIGN KEY (external_user_id)
        REFERENCES profile_cache(external_user_id)
        ON DELETE CASCADE
);

CREATE TABLE processed_events(
    event_id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    event_type VARCHAR(50) NOT NULL,
    processed_at DATETIME NOT NULL,
    result VARCHAR(20) NOT NULL
);