-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS "address"
(
    "id"             uuid PRIMARY KEY,
    "name"           text      NOT NULL,
    "address_line_1" text      NOT NULL,
    "address_line_2" text      NOT NULL,
    "neighborhood"   text      NOT NULL,
    "city"           text      NOT NULL,
    "state"          text      NOT NULL,
    "postal_code"    text      NOT NULL,
    "location"       geography(POINT, 4326),
    "country"        text      NOT NULL,
    "created_at"     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at"     timestamptz
);

CREATE TABLE "customer_address" (
    "customer_id" uuid NOT NULL REFERENCES customer(id),
    "address_id"  uuid NOT NULL REFERENCES address(id),
    "is_default"  boolean NOT NULL DEFAULT false,
    PRIMARY KEY ("customer_id", "address_id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "address" CASCADE;
-- +goose StatementEnd
