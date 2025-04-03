CREATE TABLE IF NOT EXISTS "customer" (
                                          "id" text PRIMARY KEY,
                                          "name" text NOT NULL,
                                          "last_name" text NOT NULL,
                                          "cpf" text UNIQUE NOT NULL,
                                          "email" text UNIQUE NOT NULL,
                                          "phone" text NOT NULL,
                                          "created_at" timestamp NOT NULL,
                                          "updated_at" timestamp,
                                          "deleted_at" timestamp
);