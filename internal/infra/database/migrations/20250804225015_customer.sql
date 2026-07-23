-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "customer" (
    "id" UUID PRIMARY KEY,
    "external_id" text UNIQUE NOT NULL,
    "name" text NOT NULL,
    "last_name" text NOT NULL,
    "cpf" text UNIQUE,
    "email" text UNIQUE NOT NULL,
    "phone" text NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "customer" CASCADE;
-- +goose StatementEnd
