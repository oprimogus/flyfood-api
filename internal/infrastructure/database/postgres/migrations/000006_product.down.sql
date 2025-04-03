ALTER TABLE "product" DROP CONSTRAINT "product_store_id_fkey";

DROP TABLE IF EXISTS "product";

DROP TYPE IF EXISTS "ProductType" CASCADE;