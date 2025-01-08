CREATE TYPE "StoreType" AS ENUM (
    'RESTAURANT',
    'PHARMACY',
    'TOBBACO',
    'MARKET',
    'CONVENIENCE',
    'PUB'
    );

CREATE TYPE "PaymentMethod" AS ENUM (
    'CREDIT',
    'DEBIT',
    'CASH',
    'PIX',
    'BTC'
    );


CREATE TABLE IF NOT EXISTS "store" (
     "id" uuid PRIMARY KEY,
     "owner_id" BIGSERIAL NOT NULL,
     "cnpj" varchar UNIQUE NOT NULL,
     "name" varchar UNIQUE NOT NULL,
     "description" varchar UNIQUE NOT NULL,
     "active" bool NOT NULL,
     "phone" varchar UNIQUE NOT NULL,
     "score" int NOT NULL,
     "is_open" boolean NOT NULL,
     "type" "StoreType" NOT NULL,
     "profile_image" varchar,
     "header_image" varchar,
     "address_line_1" varchar NOT NULL,
     "address_line_2" varchar NOT NULL,
     "neighborhood" varchar NOT NULL,
     "city" varchar NOT NULL,
     "state" varchar NOT NULL,
     "postal_code" varchar NOT NULL,
     "latitude" varchar,
     "longitude" varchar,
     "country" varchar NOT NULL,
     "created_at" timestamp NOT NULL,
     "updated_at" timestamp,
     "deleted_at" timestamp
);

CREATE TABLE IF NOT EXISTS "store_payment_method" (
    "id" uuid,
    "payment_method" "PaymentMethod" NOT NULL
);

CREATE TABLE IF NOT EXISTS "store_business_hour" (
    "id" uuid,
    "weekday" int NOT NULL,
    "open_hour" varchar NOT NULL,
    "closing_hour" varchar NOT NULL
);

ALTER TABLE "store" ADD FOREIGN KEY ("owner_id") REFERENCES "owner" ("id");

ALTER TABLE "store_payment_method" ADD FOREIGN KEY ("id") REFERENCES "store" ("id");

ALTER TABLE "store_payment_method" ADD CONSTRAINT "unique_store_payment_method" UNIQUE ("id", "payment_method");

ALTER TABLE "store_business_hour" ADD FOREIGN KEY ("id") REFERENCES "store" ("id");

ALTER TABLE "store_business_hour" ADD CONSTRAINT "unique_store_business_hour" UNIQUE ("id", "weekday", "open_hour", "closing_hour");
