-- +goose Up
-- +goose StatementBegin
BEGIN;

CREATE TABLE teams
(
    guid        UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    description TEXT         NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ  NULL
);

CREATE UNIQUE INDEX idx_teams_name_active ON teams (name) WHERE deleted_at IS NULL;

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TABLE IF EXISTS teams;

COMMIT;
-- +goose StatementEnd
