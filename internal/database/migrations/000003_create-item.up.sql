-- Migration 000003
CREATE TYPE "ItemType" AS ENUM (
  'FOOD',
  'WATERGALLON'
);

CREATE TABLE "item" (
  "id" bigserial UNIQUE PRIMARY KEY,
  "store_id" uuid NOT NULL,
  "type" "ItemType" NOT NULL,
  "name" varchar(25) NOT NULL,
  "description" varchar(50) NOT NULL,
  "score" int NOT NULL,
  "active" bool NOT NULL,
  "discount_active" bool NOT NULL,
  "image" varchar,
  "detail" JSONB,
  "price" int NOT NULL,
  "discount_price" int NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

ALTER TABLE "item" ADD FOREIGN KEY ("store_id") REFERENCES "store" ("id");