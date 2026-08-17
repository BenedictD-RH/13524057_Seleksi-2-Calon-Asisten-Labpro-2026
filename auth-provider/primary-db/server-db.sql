CREATE TABLE users(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    name VARCHAR(20) NOT NULL,
    email VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(64) NOT NULL,
    status ENUM('Active', 'Inactive') NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE groups_(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    name VARCHAR(20) NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE user_groups(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    user_id BINARY(16),
    group_id BINARY(16),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_ug_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_ug_group
        FOREIGN KEY (group_id)
        REFERENCES groups_(id)
        ON DELETE CASCADE,
    UNIQUE KEY uq_user_groups (user_id, group_id) 
);

CREATE TABLE applications(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    name VARCHAR(20) NOT NULL,
    client_id VARCHAR(20) NOT NULL UNIQUE,
    client_secret_hash VARCHAR(60),
    status ENUM('Active', 'Inactive') NOT NULL,
    launch_url TEXT,
    logout_notification_url TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE application_redirect_uris(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    application_id BINARY(16) NOT NULL,
    redirect_uri TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_aru_application
        FOREIGN KEY (application_id)
        REFERENCES applications(id)
        ON DELETE CASCADE
);

CREATE TABLE application_group_policies(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    application_id BINARY(16) NOT NULL,
    group_id BINARY(16) NOT NULL,
    effect ENUM('Allow', 'Blocked') NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_agp_application
        FOREIGN KEY (application_id)
        REFERENCES applications(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_agp_group
        FOREIGN KEY (group_id)
        REFERENCES groups_(id)
        ON DELETE CASCADE,
    UNIQUE KEY uq_application_group (application_id, group_id) 
);


CREATE TABLE sso_sessions(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    user_id BINARY(16) NOT NULL,
    session_token_hash VARCHAR(64) NOT NULL,
    status ENUM('Active', 'Expired', 'Revoked') NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at DATETIME NOT NULL,
    last_activity_at DATETIME,
    revoked_at DATETIME,
    revoked_reason VARCHAR(100),
    ip_address VARCHAR(50),
    user_agent TEXT,
    CONSTRAINT fk_sso_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE authorization_codes(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    code_hash VARCHAR(64) NOT NULL,
    code_challenge VARCHAR(64) NOT NULL,
    user_id BINARY(16) NOT NULL,
    application_id BINARY(16) NOT NULL,
    sso_session_id BINARY(16) NOT NULL,
    redirect_uri TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at DATETIME NOT NULL,
    used_at DATETIME,
    CONSTRAINT fk_auth_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_auth_application
        FOREIGN KEY (application_id)
        REFERENCES applications(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_auth_sso
        FOREIGN KEY (sso_session_id)
        REFERENCES sso_sessions(id)
        ON DELETE CASCADE
);

CREATE TABLE access_tokens(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    token_hash VARCHAR(64) NOT NULL,
    user_id BINARY(16) NOT NULL,
    application_id BINARY(16) NOT NULL,
    sso_session_id BINARY(16) NOT NULL,
    scopes TEXT,
    status ENUM('Active', 'Expired', 'Revoked') NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at DATETIME NOT NULL,
    revoked_at DATETIME,
    CONSTRAINT fk_accs_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_accs_application
        FOREIGN KEY (application_id)
        REFERENCES applications(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_accs_sso
        FOREIGN KEY (sso_session_id)
        REFERENCES sso_sessions(id)
        ON DELETE CASCADE
);

CREATE TABLE audit_logs(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    event_type VARCHAR(50) NOT NULL,
    actor_id BINARY(16),
    user_id BINARY(16),
    application_id BINARY(16),
    session_id BINARY(16),
    result VARCHAR(50) NOT NULL,
    metadata JSON,
    ip_address VARCHAR(50),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_audit_actor
        FOREIGN KEY (actor_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_audit_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_audit_application
        FOREIGN KEY (application_id)
        REFERENCES applications(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_audit_sso
        FOREIGN KEY (session_id)
        REFERENCES sso_sessions(id)
        ON DELETE CASCADE
)