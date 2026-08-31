-- +goose Up
-- +goose StatementBegin
BEGIN;

CREATE TABLE workspaces
(
    guid       UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    team_guid  UUID         NOT NULL REFERENCES teams (guid) ON DELETE CASCADE,
    name       VARCHAR(255) NOT NULL,
    version    VARCHAR(50)  NOT NULL DEFAULT '1.0.0',
    data       JSONB        NOT NULL DEFAULT '{}',
    config     JSONB        NOT NULL DEFAULT '{}',
    envs       JSONB        NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ  NULL
);

CREATE INDEX idx_workspaces_team_guid ON workspaces (team_guid);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TABLE IF EXISTS workspaces;

COMMIT;
-- +goose StatementEnd
