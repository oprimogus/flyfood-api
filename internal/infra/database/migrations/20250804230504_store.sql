-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "store_type_domain" (
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO "store_type_domain" (name) VALUES
    ('RESTAURANT'),
    ('PHARMACY'),
    ('TOBACCO'),
    ('MARKET'),
    ('CONVENIENCE'),
    ('PUB');

CREATE TABLE IF NOT EXISTS "payment_method_domain" (
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO "payment_method_domain" (name) VALUES
    ('CREDIT'),
    ('DEBIT'),
    ('CASH'),
    ('PIX'),
    ('BTC');

CREATE TABLE IF NOT EXISTS "store" (
    "id"            uuid PRIMARY KEY,
    "owner_id"      uuid NOT NULL REFERENCES owner(id),
    "address_id"    uuid NOT NULL REFERENCES address(id),
    "cnpj"          text UNIQUE NOT NULL,
    "name"          text UNIQUE NOT NULL,
    "description"   text,
    "active"        bool NOT NULL DEFAULT false,
    "phone"         text UNIQUE NOT NULL,
    "score"         int NOT NULL DEFAULT 0,
    "is_open"       boolean NOT NULL DEFAULT false,
    "type"          int NOT NULL REFERENCES store_type_domain(id),
    "profile_image" text,
    "header_image"  text,
    "created_at"    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at"    timestamptz
);

CREATE TABLE IF NOT EXISTS "store_payment_method" (
    "store_id"       uuid NOT NULL REFERENCES store(id),
    "payment_method" int NOT NULL REFERENCES payment_method_domain(id),
    PRIMARY KEY ("store_id", "payment_method")
);

CREATE TABLE IF NOT EXISTS "store_business_hour" (
    "store_id"     uuid NOT NULL REFERENCES store(id),
    "weekday"      smallint NOT NULL CHECK (weekday BETWEEN 0 AND 6),
    "open_hour"    smallint NOT NULL,
    "closing_hour" smallint NOT NULL,
    PRIMARY KEY ("store_id", "weekday")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "store_business_hour";
DROP TABLE IF EXISTS "store_payment_method";
DROP TABLE IF EXISTS "store";
DROP TABLE IF EXISTS "store_type_domain";
DROP TABLE IF EXISTS "payment_method_domain";

-- +goose StatementEnd