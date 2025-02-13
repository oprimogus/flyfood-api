ALTER TABLE "store_business_hour" DROP CONSTRAINT "unique_store_business_hour";
ALTER TABLE "store_business_hour" DROP CONSTRAINT "store_business_hour_id_fkey";
ALTER TABLE "store_payment_method" DROP CONSTRAINT "unique_store_payment_method";
ALTER TABLE "store_payment_method" DROP CONSTRAINT "store_payment_method_id_fkey";
ALTER TABLE "store" DROP CONSTRAINT "store_owner_id_fkey";

DROP TABLE IF EXISTS "store_business_hour";
DROP TABLE IF EXISTS "store_payment_method";
DROP TABLE IF EXISTS "store";

DROP TYPE IF EXISTS "StoreType" CASCADE;
DROP TYPE IF EXISTS "StoreStatus" CASCADE;
DROP TYPE IF EXISTS "PaymentMethod" CASCADE;