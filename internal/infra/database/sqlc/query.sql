-- name: FindCustomerByID :one
SELECT c.id, c.external_id, c.name, c.last_name, c.cpf, c.email, c.phone
FROM customer c
WHERE c.id = $1
LIMIT 1;

-- name: FindCustomerByExternalID :one
SELECT c.id, c.external_id, c.name, c.last_name, c.cpf, c.email, c.phone
FROM customer c
WHERE c.external_id = $1
LIMIT 1;

-- name: SaveCustomer :exec
INSERT INTO customer (id, external_id, name, last_name, cpf, email, phone, created_at, updated_at)
VALUES (
    @id, @external_id, @name, @last_name, @cpf, @email, @phone, 
    NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC'
)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    last_name = EXCLUDED.last_name,
    cpf = EXCLUDED.cpf,
    email = EXCLUDED.email,
    phone = EXCLUDED.phone,
    updated_at = NOW() AT TIME ZONE 'UTC';

-- name: FindOwnerByID :one
SELECT o.id, o.signature_active, o.created_at, o.updated_at, o.deleted_at
FROM owner o
WHERE o.id = $1
LIMIT 1;

-- name: FindOwnerByCustomerExternalID :one
SELECT o.id, o.signature_active, o.created_at, o.updated_at, o.deleted_at
FROM owner o
JOIN customer c ON o.id = c.id
WHERE c.external_id = $1
LIMIT 1;

-- name: SaveOwner :exec
INSERT INTO owner (id, signature_active, created_at, updated_at)
VALUES (
    @id, @signature_active, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC'
)
ON CONFLICT (id) DO UPDATE SET
    signature_active = EXCLUDED.signature_active,
    updated_at = NOW() AT TIME ZONE 'UTC';

-- name: FindStoreByID :one
SELECT
    s.id, s.owner_id, s.cnpj, s.name, s.description, s.active, s.phone, s.score, s.is_open,
    s.type, s.profile_image, s.header_image, s.created_at, s.updated_at, s.deleted_at,
    a.address_line_1, a.address_line_2, a.neighborhood, a.city, a.state, a.postal_code, a.country, ST_X(a.location::geometry)::float8 AS latitude,
    ST_Y(a.location::geometry)::float8 AS longitude
FROM store s
LEFT JOIN address a               ON a.id = s.address_id
LEFT JOIN product p               ON p.store_id = s.id
WHERE s.id = $1
GROUP BY
    s.id, s.owner_id, s.cnpj, s.name, s.description, s.active, s.phone, s.score, s.is_open,
    s.type, s.profile_image, s.header_image, s.created_at, s.updated_at, s.deleted_at,
    a.address_line_1, a.address_line_2, a.neighborhood, a.city, a.state, a.postal_code, a.country, latitude, longitude;

-- name: FindStoreBusinessHoursByStoreId :many
SELECT bh.weekday, bh.open_hour, bh.closing_hour
FROM store_business_hour bh
WHERE bh.store_id = $1
ORDER BY bh.weekday;

-- name: FindStorePaymentMethodsByStoreId :many
SELECT pm.payment_method
FROM store_payment_method pm
WHERE pm.store_id = $1
ORDER BY pm.payment_method;

-- name: SaveStore :exec
INSERT INTO store (id, owner_id, address_id, cnpj, name, description, active, is_open, phone, 
    score, type, created_at, updated_at, deleted_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC', NULL)
ON CONFLICT (id) DO UPDATE
    SET
        owner_id = EXCLUDED.owner_id,
        address_id = EXCLUDED.address_id,
        cnpj = EXCLUDED.cnpj,
        name = EXCLUDED.name,
        description = EXCLUDED.description,
        active = EXCLUDED.active,
        is_open = EXCLUDED.is_open,
        phone = EXCLUDED.phone,
        score = EXCLUDED.score,
        type = EXCLUDED.type,
        updated_at = NOW() AT TIME ZONE 'UTC';

-- name: UpsertBusinessHour :exec
INSERT INTO store_business_hour (store_id, weekday, open_hour, closing_hour)
VALUES ($1, $2, $3, $4)
ON CONFLICT (store_id, weekday) DO UPDATE SET
    open_hour    = EXCLUDED.open_hour,
    closing_hour = EXCLUDED.closing_hour;

-- name: DeleteBusinessHoursNotIn :exec
DELETE FROM store_business_hour
WHERE store_id = $1
  AND weekday != ALL(@business_hours::smallint[]);

-- name: UpsertPaymentMethod :exec
INSERT INTO store_payment_method (store_id, payment_method)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: DeletePaymentMethodsNotIn :exec
DELETE FROM store_payment_method
WHERE store_id = $1
  AND payment_method != ALL(@payment_methods::int[]);

-- name: FindAddressesByExternalCustomerID :many
SELECT a.id, a.name, ca.is_default, a.address_line_1, a.address_line_2, a.neighborhood, a.city, 
    a.state, a.postal_code, a.country, ST_Y(a.location::geometry)::float8 AS latitude,
    ST_X(a.location::geometry)::float8 AS longitude
FROM address a
LEFT JOIN customer_address ca ON a.id = ca.address_id
LEFT JOIN customer c ON ca.customer_id = c.id
WHERE c.external_id = $1;

-- name: SaveAddress :exec
INSERT INTO address (id, name, address_line_1, address_line_2, neighborhood, 
    city, state, postal_code, country, location, created_at, updated_at, deleted_at)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9,
    ST_SetSRID(ST_MakePoint(@latitude::float8, @longitude::float8), 4326)::geography,
    NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC', NULL
)
ON CONFLICT (id) DO UPDATE
    SET
        name          = EXCLUDED.name,
        address_line_1 = EXCLUDED.address_line_1,
        address_line_2 = EXCLUDED.address_line_2,
        neighborhood  = EXCLUDED.neighborhood,
        city          = EXCLUDED.city,
        state         = EXCLUDED.state,
        postal_code   = EXCLUDED.postal_code,
        country       = EXCLUDED.country,
        location      = EXCLUDED.location,
        updated_at    = NOW() AT TIME ZONE 'UTC';

-- name: SaveCustomerAddress :exec
INSERT INTO customer_address (customer_id, address_id, is_default)
VALUES ($1, $2, $3);

-- name: DeleteAddress :one
DELETE FROM address
WHERE id = $1
RETURNING id;

-- name: DeleteCustomerAddress :one
DELETE FROM customer_address
WHERE customer_id = $1 AND address_id = $2
RETURNING customer_id;

