ALTER TABLE "address" DROP CONSTRAINT "address_customer_id_fkey";

DROP TABLE IF EXISTS "address" CASCADE;