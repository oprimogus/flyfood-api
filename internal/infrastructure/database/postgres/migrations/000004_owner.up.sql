CREATE TABLE IF NOT EXISTS "owner" (
                                       "id" text PRIMARY KEY,
                                       "signature_active" bool NOT NULL,
                                       "created_at" timestamp NOT NULL,
                                       "updated_at" timestamp,
                                       "deleted_at" timestamp
);

ALTER TABLE "owner" ADD FOREIGN KEY ("id") REFERENCES "customer" ("id");