-- MIGRATION 2
CREATE TABLE IF NOT EXISTS "customer" (
                                          "id" text PRIMARY KEY,
                                          "name" text NOT NULL,
                                          "last_name" text NOT NULL,
                                          "cpf" text UNIQUE NOT NULL,
                                          "email" text UNIQUE NOT NULL,
                                          "phone" text NOT NULL,
                                          "created_at" timestamp NOT NULL,
                                          "updated_at" timestamp,
                                          "deleted_at" timestamp
);
-- MIGRATION 3
CREATE TABLE IF NOT EXISTS "address" (
                                         "id" BIGSERIAL PRIMARY KEY,
                                         "customer_id" text NOT NULL,
                                         "name" text NOT NULL,
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

ALTER TABLE "address" ADD FOREIGN KEY ("customer_id") REFERENCES "customer" ("id");
-- MIGRATION 4
CREATE TABLE IF NOT EXISTS "owner" (
                                       "id" text PRIMARY KEY,
                                       "signature_active" bool NOT NULL,
                                       "created_at" timestamp NOT NULL,
                                       "updated_at" timestamp,
                                       "deleted_at" timestamp
);

ALTER TABLE "owner" ADD FOREIGN KEY ("id") REFERENCES "customer" ("id");
-- MIGRATION 5
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

-- MIGRATION 6
CREATE TYPE "ProductType" AS ENUM (
    'FOOD',
    'WATER'
    );

CREATE TABLE "product" (
                           "id" uuid PRIMARY KEY,
                           "store_id" uuid NOT NULL,
                           "sku" text,
                           "active_for_sale" bool NOT NULL,
                           "promo_active" bool NOT NULL,
                           "type" "ProductType" NOT NULL,
                           "tag" text NOT NULL,
                           "name" text NOT NULL,
                           "description" text NOT NULL,
                           "stock_quantity" int NOT NULL,
                           "score" int NOT NULL,
                           "image_url" text,
                           "details" jsonb,
                           "price" int NOT NULL,
                           "promotional_price" int,
                           "created_at" timestamp NOT NULL,
                           "updated_at" timestamp,
                           "deleted_at" timestamp
);

ALTER TABLE "product" ADD FOREIGN KEY ("store_id") REFERENCES "store" ("id");
-- MIGRATION 7
CREATE TYPE "OrderStatus" AS ENUM (
    'CREATED',
    'CANCELLED',
    'VERIFIED_BY_CUSTOMER',
    'VERIFIED_BY_STORE',
    'IN_PROCESS',
    'DISPATCHED',
    'DELIVERED',
    'CHARGEBACK',
    'FINISHED'
    );

CREATE TABLE IF NOT EXISTS "order" (
                                       "id" uuid PRIMARY KEY,
                                       "store_id" uuid NOT NULL,
                                       "customer_id" text NOT NULL,
                                       "address_id" int NOT NULL,
                                       "status" "OrderStatus" NOT NULL,
                                       "shipping_amount" int NOT NULL,
                                       "created_at" timestamp NOT NULL,
                                       "updated_at" timestamp,
                                       "deleted_at" timestamp
);

ALTER TABLE "order" ADD FOREIGN KEY ("store_id") REFERENCES "store" ("id");
ALTER TABLE "order" ADD FOREIGN KEY ("customer_id") REFERENCES "customer" ("id");
ALTER TABLE "order" ADD FOREIGN KEY ("address_id") REFERENCES "address" ("id");

CREATE TABLE IF NOT EXISTS "order_item" (
                                            "id" uuid PRIMARY KEY,
                                            "order_id" uuid,
                                            "product_id" uuid,
                                            "quantity" int NOT NULL,
                                            "amount" int NOT NULL,
                                            "details" text
);

ALTER TABLE "order_item" ADD FOREIGN KEY ("order_id") REFERENCES "order" ("id");
ALTER TABLE "order_item" ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id");