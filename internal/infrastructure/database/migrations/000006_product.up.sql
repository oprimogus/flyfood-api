CREATE TYPE "ProductType" AS ENUM (
    'FOOD',
    'WATER'
    );

CREATE TABLE "product" (
    "id" uuid PRIMARY KEY,
    "store_id" uuid,
    "sku" varchar,
    "active_for_sale" bool NOT NULL,
    "promo_active" bool NOT NULL,
    "type" "ProductType" NOT NULL,
    "name" varchar NOT NULL,
    "description" varchar NOT NULL,
    "stock_quantity" int NOT NULL,
    "score" int NOT NULL,
    "image_url" varchar,
    "details" jsonb,
    "price" int NOT NULL,
    "promotional_price" int,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp,
    "deleted_at" timestamp
);

ALTER TABLE "product" ADD FOREIGN KEY ("store_id") REFERENCES "store" ("id");