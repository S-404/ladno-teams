-- +goose Up
-- +goose StatementBegin
BEGIN;

CREATE TYPE role_enum AS ENUM ('MAINTAINER', 'DEVELOPER', 'GUEST');

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TYPE IF EXISTS role_enum;

COMMIT;
-- +goose StatementEnd
