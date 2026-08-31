-- +goose Up
-- +goose StatementBegin
BEGIN;

CREATE TABLE users
(
    guid       UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    login      VARCHAR(100) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    is_admin   BOOLEAN      NOT NULL DEFAULT FALSE,
    is_blocked BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ  NULL
);

CREATE UNIQUE INDEX idx_users_login_active ON users (login) WHERE deleted_at IS NULL;

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TABLE IF EXISTS users;

COMMIT;
-- +goose StatementEnd
