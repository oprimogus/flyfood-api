-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "owner" (
   "id" UUID PRIMARY KEY,
   "signature_active" bool NOT NULL,
   "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
   "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
   "deleted_at" timestamptz
);
ALTER TABLE "owner" ADD FOREIGN KEY ("id") REFERENCES "customer" ("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "owner" DROP CONSTRAINT "owner_id_fkey";
DROP TABLE IF EXISTS "owner" CASCADE;
-- +goose StatementEnd
