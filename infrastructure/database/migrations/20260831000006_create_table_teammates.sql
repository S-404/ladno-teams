-- +goose Up
-- +goose StatementBegin
BEGIN;

CREATE TABLE teammates
(
    user_guid  UUID NOT NULL REFERENCES users (guid) ON DELETE CASCADE,
    team_guid  UUID NOT NULL REFERENCES teams (guid) ON DELETE CASCADE,
    is_leader  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_guid, team_guid)
);

CREATE INDEX idx_teammates_team_guid ON teammates (team_guid);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TABLE IF EXISTS teammates;

COMMIT;
-- +goose StatementEnd
