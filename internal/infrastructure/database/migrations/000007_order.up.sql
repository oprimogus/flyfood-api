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