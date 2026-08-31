-- +goose Up
-- +goose StatementBegin
BEGIN;

CREATE TABLE invites
(
    guid       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_guid  UUID NOT NULL REFERENCES teams (guid) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expired_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_invites_team_guid ON invites (team_guid);
CREATE INDEX idx_invites_expired_at ON invites (expired_at);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TABLE IF EXISTS invites;

COMMIT;
-- +goose StatementEnd
