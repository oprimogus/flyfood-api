CREATE TABLE IF NOT EXISTS "address" (
   "id" BIGSERIAL PRIMARY KEY,
   "customer_id" uuid,
   "name" varchar NOT NULL,
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

ALTER TABLE "address" ADD FOREIGN KEY ("customer_id") REFERENCES "customer" ("id");