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
                                       "owner_id" text NOT NULL,
                                       "cnpj" text UNIQUE NOT NULL,
                                       "name" text UNIQUE NOT NULL,
                                       "description" text UNIQUE NOT NULL,
                                       "active" bool NOT NULL,
                                       "phone" text UNIQUE NOT NULL,
                                       "score" int NOT NULL,
                                       "is_open" boolean NOT NULL,
                                       "type" "StoreType" NOT NULL,
                                       "profile_image" text,
                                       "header_image" text,
                                       "address_line_1" text NOT NULL,
                                       "address_line_2" text NOT NULL,
                                       "neighborhood" text NOT NULL,
                                       "city" text NOT NULL,
                                       "state" text NOT NULL,
                                       "postal_code" text NOT NULL,
                                       "latitude" text,
                                       "longitude" text,
                                       "country" text NOT NULL,
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
                                                     "open_hour" text NOT NULL,
                                                     "closing_hour" text NOT NULL
);

ALTER TABLE "store" ADD FOREIGN KEY ("owner_id") REFERENCES "owner" ("id");

ALTER TABLE "store_payment_method" ADD FOREIGN KEY ("id") REFERENCES "store" ("id");

ALTER TABLE "store_payment_method" ADD CONSTRAINT "unique_store_payment_method" UNIQUE ("id", "payment_method");

ALTER TABLE "store_business_hour" ADD FOREIGN KEY ("id") REFERENCES "store" ("id");

ALTER TABLE "store_business_hour" ADD CONSTRAINT "unique_store_business_hour" UNIQUE ("id", "weekday", "open_hour", "closing_hour");
