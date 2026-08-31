-- +goose Up
-- +goose StatementBegin
BEGIN;

CREATE TABLE profiles
(
    user_guid UUID PRIMARY KEY REFERENCES users (guid) ON DELETE CASCADE,
    name      VARCHAR(255) NOT NULL,
    about     TEXT         NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TABLE IF EXISTS profiles;

COMMIT;
-- +goose StatementEnd
