-- +goose Up
-- +goose StatementBegin
BEGIN;

CREATE TABLE workspace_roles
(
    workspace_guid UUID NOT NULL REFERENCES workspaces (guid) ON DELETE CASCADE,
    teammate_guid  UUID NOT NULL REFERENCES users (guid) ON DELETE CASCADE,
    role           role_enum NOT NULL DEFAULT 'GUEST',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (workspace_guid, teammate_guid)
);

CREATE INDEX idx_workspace_roles_workspace_guid ON workspace_roles (workspace_guid);
CREATE INDEX idx_workspace_roles_teammate_guid ON workspace_roles (teammate_guid);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TABLE IF EXISTS workspace_roles;

COMMIT;
-- +goose StatementEnd
