-- +goose Up
-- +goose StatementBegin
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
   "id" UUID PRIMARY KEY,
   "store_id" UUID NOT NULL,
   "customer_id" UUID NOT NULL,
   "address_id" UUID NOT NULL,
   "status" "OrderStatus" NOT NULL,
   "shipping_amount" int NOT NULL,
   "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
   "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
   "deleted_at" timestamptz
);

ALTER TABLE "order" ADD FOREIGN KEY ("store_id") REFERENCES "store" ("id");
ALTER TABLE "order" ADD FOREIGN KEY ("customer_id") REFERENCES "customer" ("id");
ALTER TABLE "order" ADD FOREIGN KEY ("address_id") REFERENCES "address" ("id");

CREATE TABLE IF NOT EXISTS "order_item" (
    "id" UUID PRIMARY KEY,
    "order_id" UUID,
    "product_id" UUID,
    "quantity" int NOT NULL,
    "amount" int NOT NULL,
    "details" text
);

ALTER TABLE "order_item" ADD FOREIGN KEY ("order_id") REFERENCES "order" ("id");
ALTER TABLE "order_item" ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS "order_item" DROP CONSTRAINT "order_item_order_id_fkey";
ALTER TABLE IF EXISTS "order_item" DROP CONSTRAINT "order_item_product_id_fkey";

DROP TABLE IF EXISTS "order_item";

ALTER TABLE IF EXISTS "order" DROP CONSTRAINT "order_store_id_fkey" CASCADE;
ALTER TABLE IF EXISTS "order" DROP CONSTRAINT "order_customer_id_fkey" CASCADE;
ALTER TABLE IF EXISTS "order" DROP CONSTRAINT "order_address_id_fkey" CASCADE;

DROP TABLE IF EXISTS "order" CASCADE;

DROP TYPE IF EXISTS "OrderStatus" CASCADE;
-- +goose StatementEnd
