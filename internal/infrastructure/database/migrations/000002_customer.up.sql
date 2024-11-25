CREATE TABLE IF NOT EXISTS "customer" (
    "id" BIGSERIAL PRIMARY KEY,
    "name" varchar NOT NULL,
    "last_name" varchar NOT NULL,
    "cpf" varchar UNIQUE NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "phone" varchar NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp,
    "deleted_at" timestamp
);