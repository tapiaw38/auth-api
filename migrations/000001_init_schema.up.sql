CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(32) PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255),
    picture VARCHAR(255),
    address VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    verified_email BOOLEAN NOT NULL DEFAULT FALSE,
    verified_email_token VARCHAR(255) UNIQUE NOT NULL DEFAULT '',
    verified_email_token_expiry TIMESTAMP NOT NULL DEFAULT NOW(),
    password_reset_token VARCHAR(255) UNIQUE NOT NULL DEFAULT '',
    password_reset_token_expiry TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS roles (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL 
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id VARCHAR(32) REFERENCES users(id) ON DELETE CASCADE,
    role_id VARCHAR(32) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT user_roles_pkey PRIMARY KEY (user_id, role_id)
);