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