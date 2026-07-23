-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS unaccent;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd
