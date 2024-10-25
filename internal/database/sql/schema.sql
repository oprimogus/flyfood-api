-- Migration 000001
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Migration 000002
CREATE TYPE "ShopType" AS ENUM (
  'restaurant',
  'pharmacy',
  'tobbaco',
  'market',
  'convenience',
  'pub'
);

CREATE TYPE "PaymentMethodEnum" AS ENUM (
  'credit',
  'debit',
  'pix',
  'cash',
  'bitcoin'
);

CREATE TABLE "store" (
  "id" uuid UNIQUE PRIMARY KEY,
  "cpf_cnpj" varchar UNIQUE NOT NULL,
  "owner_id" uuid NOT NULL,
  "name" varchar(25) UNIQUE NOT NULL,
  "active" bool NOT NULL,
  "phone" varchar(18) UNIQUE NOT NULL,
  "score" int NOT NULL,
  "type" "ShopType" NOT NULL,
  "address_line_1" varchar(40) NOT NULL,
  "address_line_2" varchar(20) NOT NULL,
  "neighborhood" varchar(25) NOT NULL,
  "city" varchar(25) NOT NULL,
  "state" varchar(15) NOT NULL,
  "postal_code" varchar(15) NOT NULL,
  "latitude" varchar,
  "longitude" varchar,
  "country" varchar(15) NOT NULL,
  "profile_image" varchar UNIQUE,
  "header_image" varchar UNIQUE,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "business_hour" (
  "store_id" uuid NOT NULL,
  "week_day" int NOT NULL CHECK ("week_day" >= 0 AND "week_day" <= 6),
  "opening_time" TIME NOT NULL,
  "closing_time" TIME NOT NULL,
  "timezone" VARCHAR NOT NULL,
  CONSTRAINT unique_business_hour UNIQUE ("store_id", "week_day", "opening_time", "closing_time")
);

CREATE TABLE "payment_method" (
  "id" integer PRIMARY KEY,
  "method" "PaymentMethodEnum"
);

CREATE TABLE "store_payment_method" (
  "store_id" uuid PRIMARY KEY,
  "payment_method_id" integer
);

CREATE TABLE "address" (
  "id" integer PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "address_line_1" varchar NOT NULL,
  "address_line_2" varchar NOT NULL,
  "neighborhood" varchar NOT NULL,
  "city" varchar NOT NULL,
  "state" varchar NOT NULL,
  "postal_code" varchar NOT NULL,
  "latitude" varchar,
  "longitude" varchar NOT NULL,
  "country" varchar NOT NULL,
  "created_at" timestamp NOT NULL
);

COMMENT ON COLUMN "store"."owner_id" IS 'user ID of store owner';

ALTER TABLE "business_hour" ADD FOREIGN KEY ("store_id") REFERENCES "store" ("id");

ALTER TABLE "store_payment_method" ADD FOREIGN KEY ("store_id") REFERENCES "store" ("id");

ALTER TABLE "store_payment_method" ADD FOREIGN KEY ("payment_method_id") REFERENCES "payment_method" ("id");

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