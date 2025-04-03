ALTER TABLE IF EXISTS "order_item" DROP CONSTRAINT "order_item_order_id_fkey";
ALTER TABLE IF EXISTS "order_item" DROP CONSTRAINT "order_item_product_id_fkey";

DROP TABLE IF EXISTS "order_item";

ALTER TABLE IF EXISTS "order" DROP CONSTRAINT "order_store_id_fkey" CASCADE;
ALTER TABLE IF EXISTS "order" DROP CONSTRAINT "order_customer_id_fkey" CASCADE;
ALTER TABLE IF EXISTS "order" DROP CONSTRAINT "order_address_id_fkey" CASCADE;

DROP TABLE IF EXISTS "order" CASCADE;

DROP TYPE IF EXISTS "OrderStatus" CASCADE;