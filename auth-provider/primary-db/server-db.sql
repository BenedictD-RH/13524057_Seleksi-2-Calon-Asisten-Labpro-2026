CREATE TABLE users(
    id BINARY(16) PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    email VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(50) NOT NULL,
    status ENUM('Active', 'Inactive') NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE groups_(
    id BINARY(16) PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE user_groups(
    id BINARY(16) PRIMARY KEY,
    user_id BINARY(16),
    group_id BINARY(16),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id),
    CONSTRAINT fk_group
        FOREIGN KEY (group_id)
        REFERENCES groups_(id)
)
