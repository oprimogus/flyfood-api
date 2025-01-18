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