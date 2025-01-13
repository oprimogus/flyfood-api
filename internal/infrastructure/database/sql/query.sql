-- noinspection SqlResolveForFile

-- name: FindCustomerByID :one
SELECT c.id, c.name, c.last_name, c.cpf, c.email, c.phone
FROM customer c
WHERE c.id = $1
LIMIT 1;

-- name: FindAddressesByCustomerID :many
SELECT a.id, a.name, a.address_line_1, address_line_2, a.neighborhood, a.city, a.state, a.postal_code, a.country, a.latitude, a.longitude
FROM address a
WHERE a.customer_id = $1
ORDER BY a.created_at DESC;

-- name: SetOwner :exec
INSERT INTO owner (id, signature_active, created_at, updated_at, deleted_at)
VALUES ($1, $2, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC', null);

-- name: SaveOwner :exec
INSERT INTO owner (id, signature_active, created_at, updated_at, deleted_at)
VALUES ($1, $2, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC', null)
ON CONFLICT (id) DO UPDATE
SET
    signature_active = excluded.signature_active,
    updated_at = NOW() AT TIME ZONE 'UTC';

-- name: RemoveOwner :exec
DELETE FROM owner
WHERE id = $1;

-- name: SaveCustomer :exec
INSERT INTO customer (id, name, last_name, cpf, email, phone, created_at, updated_at, deleted_at)
VALUES ($1, $2, $3, $4, $5, $6,
        NOW() AT TIME ZONE 'UTC',
        NOW() AT TIME ZONE 'UTC',
        NULL)
ON CONFLICT (id) DO UPDATE
    SET
        name = EXCLUDED.name,
        last_name = EXCLUDED.last_name,
        cpf = EXCLUDED.cpf,
        email = EXCLUDED.email,
        phone = EXCLUDED.phone
    WHERE customer.id = $1;

-- name: SaveCustomerAddress :exec
INSERT INTO address (customer_id, name, address_line_1, address_line_2, neighborhood, city, state,
                     postal_code, latitude, longitude, country, created_at, updated_at, deleted_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8,
        $9, $10, $11, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC', NULL)
ON CONFLICT (customer_id, address_line_1, address_line_2, postal_code) DO UPDATE
    SET
        name = EXCLUDED.name,
        neighborhood = EXCLUDED.neighborhood,
        city = EXCLUDED.city,
        state = EXCLUDED.state,
        latitude = EXCLUDED.latitude,
        longitude = EXCLUDED.longitude,
        country = EXCLUDED.country,
        updated_at = NOW() AT TIME ZONE 'UTC'
WHERE address.customer_id = $1
  AND address.address_line_1 = $3
  AND address.address_line_2 = $4
  AND address.postal_code = $8;

-- name: AddNewCustomerAddresses :batchexec
INSERT INTO address (customer_id, name, address_line_1, address_line_2, neighborhood, city, state,
postal_code, latitude, longitude, country, created_at, updated_at, deleted_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8,
        $9, $10, $11, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC', NULL);

-- name: DeleteCustomerAddresses :batchexec
DELETE FROM address
WHERE customer_id = $1
    AND name = $2
    AND address_line_1 = $3
    AND address_line_2 = $4
    AND neighborhood = $5
    AND city = $6
  And state = $7
    AND postal_code = $8
    AND country = $9;

-- name: FindStoreByID :one
WITH relevant_business_hours AS (
    SELECT *
    FROM store_business_hour bh
    WHERE id = $1
),
     relevant_payment_methods AS (
         SELECT *
         FROM store_payment_method pm
         WHERE id = $1
     )
SELECT
    s.owner_id, s.cnpj, s.name, s.description,
    s.active, s.phone, s.score, s.is_open, s.type, s.profile_image,
    s.header_image, s.address_line_1, s.address_line_2, s.neighborhood,
    s.city, s.state, s.postal_code, s.latitude, s.longitude,
    s.country, s.created_at, s.updated_at, s.deleted_at,
    COALESCE(
        ARRAY_AGG(DISTINCT rbh.weekday || ' ' || rbh.open_hour || ' ' || rbh.closing_hour) FILTER (WHERE rbh.id IS NOT NULL),
                    '{}'::text[]
    )::text[] AS business_hours,
    COALESCE(
        ARRAY_AGG(DISTINCT rpm.payment_method) FILTER (WHERE rpm.id IS NOT NULL),
                    '{}'::"PaymentMethod"[]
    )::text[] AS payment_methods
FROM store s
         LEFT JOIN relevant_business_hours rbh ON s.id = rbh.id
         LEFT JOIN relevant_payment_methods rpm ON s.id = rpm.id
WHERE s.id = $1
GROUP BY
    s.owner_id, s.cnpj, s.name, s.description,
    s.active, s.phone, s.score, s.is_open, s.type, s.profile_image,
    s.header_image, s.address_line_1, s.address_line_2, s.neighborhood,
    s.city, s.state, s.postal_code, s.latitude, s.longitude,
    s.country, s.created_at, s.updated_at, s.deleted_at;

-- name: FindBusinessHourByStoreID :many
SELECT bh.weekday, bh.open_hour, bh.closing_hour
FROM store_business_hour bh
WHERE bh.id = $1;

-- name: FindPaymentMethodsByStoreID :many
SELECT pm.payment_method FROM store_payment_method pm
WHERE id = $1;

-- name: FindProductsIDByStoreID :many
SELECT p.id from product p
WHERE p.store_id = $1;

-- name: SaveStore :exec
INSERT INTO store (id, owner_id, cnpj, name, description, active, phone, score, is_open, type,
                   profile_image, header_image, address_line_1, address_line_2,
                   neighborhood, city, state, postal_code, country, latitude, longitude,
                   created_at, updated_at, deleted_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
        $11, $12, $13, $14, $15, $16,
        $17, $18, $19, $20, $21,
        NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC', null)
ON CONFLICT (id) DO UPDATE
    SET
        cnpj = excluded.cnpj,
        name = excluded.name,
        description = excluded.description,
        active = excluded.active,
        phone = excluded.phone,
        score = excluded.score,
        is_open = excluded.is_open,
        type = excluded.type,
        profile_image = excluded.profile_image,
        header_image = excluded.header_image,
        address_line_1 = excluded.address_line_1,
        address_line_2 = excluded.address_line_2,
        neighborhood = excluded.neighborhood,
        city = excluded.city,
        state = excluded.state,
        postal_code = excluded.postal_code,
        country = excluded.country,
        latitude = excluded.latitude,
        longitude = excluded.longitude,
        updated_at = NOW() AT TIME ZONE 'UTC'
WHERE store.id = excluded.id;

-- name: IsOwner :one
SELECT EXISTS(SELECT 1 FROM owner WHERE id = $1);

-- name: FindOwnerByID :one
SELECT o.id, signature_active from owner o
WHERE o.id = $1;

-- name: SaveStoreBusinessHour :batchexec
INSERT INTO store_business_hour(id, weekday, open_hour, closing_hour)
VALUES ($1, $2, $3, $4);

-- name: DeleteStoreBusinessHour :batchexec
DELETE FROM store_business_hour bh
WHERE bh.id = $1
AND bh.weekday = $2
AND bh.open_hour = $3
AND bh.closing_hour = $4;

-- name: SaveStorePaymentMethods :batchexec
INSERT INTO store_payment_method (id, payment_method)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: DeleteStorePaymentMethods :batchexec
DELETE FROM store_payment_method
WHERE id = $1 AND payment_method = $2;

-- name: FindProductByID :one
SELECT i.store_id, i.sku, i.active_for_sale, i.promo_active, i.type, i.tag, i.name,
       i.description, i.stock_quantity, i.score, i.image_url, i.details, i.price, i.promotional_price
FROM product i
WHERE i.id = $1
LIMIT 1;

-- name: SaveProduct :exec
INSERT INTO product (id, store_id, sku, active_for_sale, promo_active, type, tag, name, description,
                     stock_quantity, score, image_url, details, price, promotional_price, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7,
        $8, $9, $10, $11, $12, $13, $14, $15,
        NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC')
ON CONFLICT (id) DO UPDATE
SET
    sku = excluded.sku,
    active_for_sale = excluded.active_for_sale,
    promo_active = excluded.promo_active,
    type = excluded.type,
    tag = excluded.tag,
    name = excluded.name,
    description = excluded.description,
    stock_quantity = excluded.stock_quantity,
    score = excluded.score,
    image_url = excluded.image_url,
    details = excluded.details,
    price = excluded.price,
    promotional_price = excluded.promotional_price,
    updated_at = NOW() AT TIME ZONE 'UTC'
WHERE product.id = $1 AND product.store_id = $2;

--                    latitude, longitude, country, created_at, updated_at)

-- -- name: FindCustomerByID :many
-- SELECT c.id, c.name, c.last_name, c.cpf, c.email, c.phone, a.address_line_1, a.address_line_2,
--        a.neighborhood, a.city, a.state, a.postal_code, a.country, o.id
-- FROM customer c
-- LEFT JOIN address a on c.id = a.customer_id
-- LEFT JOIN "order" o on c.id = o.customer_id;
--
-- -- name: SaveCustomer :exec
-- INSERT INTO customer (id, name, last_name, cpf, email, phone, created_at, updated_at, deleted_at)
--     VALUES ($1, $2, $3, $3, $5, $6,
--             NOW() AT TIME ZONE 'UTC',
--             NOW() AT TIME ZONE 'UTC',
--             NULL)
--     ON CONFLICT (id) DO UPDATE
-- SET
--     name = $2,
--     last_name = $3,
--     cpf = $4,
--     email = $5,
--     phone = $6;




-- -- name: CreateStore :exec
-- INSERT INTO store (id, cpf_cnpj, owner_id, name, active, phone, score, type,
--                    address_line_1, address_line_2, neighborhood, city, state, postal_code,
--                    latitude, longitude, country, created_at, updated_at)
-- VALUES(
--        $1,
--        $2,
--        $3,
--        $4,
--        $5,
--        $6,
--        $7,
--        $8,
--        $9,
--        $10,
--        $11,
--        $12,
--        $13,
--        $14,
--        $15,
--        $16,
--        $17,
--        NOW() AT TIME ZONE 'UTC',
--        NOW() AT TIME ZONE 'UTC');
--
-- -- name: UpdateStore :exec
-- UPDATE store
--   SET
--     name = $3,
--     phone = $4,
--     type = $5,
--     address_line_1 = $6,
--     address_line_2 = $7,
--     neighborhood = $8,
--     city = $9,
--     state = $10,
--     postal_code = $11,
--     country = $12,
--     updated_at = NOW() AT TIME ZONE 'UTC'
-- WHERE id = $1 AND owner_id = $2;
--
-- -- name: SetProfileImage :exec
-- UPDATE store
--   SET
--     profile_image = $2
-- WHERE id = $1;
--
-- -- name: SetHeaderImage :exec
-- UPDATE store
--   SET
--     header_image = $2
-- WHERE id = $1;
--
-- -- name: IsOwner :one
-- SELECT EXISTS(SELECT 1 FROM store WHERE id = $1 AND owner_id = $2);
--
-- -- name: GetStoreByID :one
-- SELECT s.id, s.name, s.phone, s.score, s.type, s.address_line_1,
-- s.address_line_2, s.neighborhood, s.city, s.state, s.country, s.profile_image, s.header_image
-- FROM store s
-- WHERE id = $1;
--
-- -- name: GetStoreBusinessHoursByID :many
-- SELECT week_day, timezone, opening_time, closing_time
-- FROM business_hour
-- WHERE store_id = $1
-- ORDER BY week_day;
--

-- name: FindStoresByFilter :many
SELECT
    s.id,
    s.name,
    s.score,
    s.is_open,
    s.type,
    s.neighborhood,
    s.latitude,
    s.longitude,
    s.profile_image
FROM store s
WHERE (CASE
           WHEN sqlc.narg('name')::text IS NOT NULL
               THEN unaccent(s.name) ILIKE '%' || unaccent(sqlc.narg('name')::text) || '%'
           ELSE true
    END)
  AND (CASE
           WHEN sqlc.narg('type')::text IS NOT NULL
               THEN s.type::text = sqlc.narg('type')::text
           ELSE true
    END)
  AND (CASE
           WHEN sqlc.narg('city')::text IS NOT NULL
               THEN s.city = sqlc.narg('city')::text
           ELSE true
    END)
  AND (CASE
           WHEN sqlc.narg('is_open')::boolean IS NOT NULL
               THEN s.is_open = sqlc.narg('is_open')::boolean
           ELSE true
    END)
ORDER BY s.score DESC, s.type
OFFSET sqlc.arg('offset_value')
LIMIT sqlc.arg('limit_items');
-- SELECT s.id, s.name, s.score, s.is_open, s.type, s.neighborhood, s.latitude, s.longitude, s.profile_image
-- FROM store s
-- WHERE 1 = 1
--     AND (COALESCE(NULLIF(@name::text, ''), s.name) IS NULL OR s.name ILIKE '%' || @name::text || '%')
--     AND (NULLIF(@type::text, '') IS NULL OR s.type::text = @type::text)
--     AND (NULLIF(@city::text, '') IS NULL OR s.city::text = @city::text)
--     AND (sqlc.narg('is_open')::bool IS NULL OR s.is_open = @is_open::bool)
-- ORDER BY s.score DESC, s.type
-- OFFSET @offset_value LIMIT @limit_items;
--
-- -- name: AddBusinessHours :batchexec
-- INSERT INTO business_hour(store_id, week_day, opening_time, closing_time, timezone)
-- VALUES ($1, $2, $3, $4, $5);
--
-- -- name: DeleteBusinessHours :batchexec
-- DELETE FROM business_hour
-- WHERE store_id = $1
--   AND week_day = $2
--   AND opening_time = $3
--   AND closing_time = $4;
--
-- -- name: FindStoreBusinessHoursByStoreId :many
-- SELECT bh.store_id, bh.week_day, bh.opening_time, bh.closing_time, bh.timezone
-- FROM business_hour bh
-- WHERE 1 = 1
--   AND bh.store_id = ANY($1::UUID[]);
--
-- -- name: CreateItem :one
-- INSERT INTO item(
--   store_id,
--   type,
--   name,
--   description,
--   score,
--   active,
--   discount_active,
--   price,
--   discount_price,
--   created_at,
--   updated_at)
-- VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,
--   NOW() AT TIME ZONE 'UTC',
--   NOW() AT TIME ZONE 'UTC')
-- RETURNING ID;
--
-- -- name: UpdateItem :exec
-- UPDATE item
--   SET
--     type = $2,
--     name = $3,
--     description = $4,
--     active = $5,
--     discount_active = $6,
--     price = $7,
--     discount_price = $8,
--     updated_at = NOW() AT TIME ZONE 'UTC'
-- WHERE id = $1;
--
-- -- name: DeleteItem :exec
-- UPDATE item
--   SET
--     deleted_at = NOW() AT TIME ZONE 'UTC'
-- WHERE id = $1;
--
-- -- name: GetItemByID :one
-- SELECT
--     i.id,
--     i.store_id,
--     i.type,
--     i.name,
--     i.description,
--     i.score,
--     i.active,
--     i.discount_active,
--     i.image,
--     i.detail,
--     i.price,
--     i.discount_price,
--     i.created_at,
--     i.updated_at,
--     i.deleted_at
-- FROM item i
-- WHERE id = $1;
--
-- -- name: GetItemByFilter :many
-- SELECT
--     i.id,
--     i.store_id,
--     i.type,
--     i.name,
--     i.score,
--     i.discount_active,
--     i.discount_price,
--     i.price,
--     i.image,
--     CASE
--         WHEN i.discount_active = true THEN i.discount_price
--         ELSE i.price
--         END AS final_price,
--     s.name AS store_name,
--     s.score AS store_score,
--     s.profile_image
-- FROM item i
--          INNER JOIN store s ON s.id = i.store_id
-- WHERE 1 = 1
--   AND s.city = @city::text
--   AND i.active = true
--   AND i.deleted_at IS NULL
--   AND (COALESCE(NULLIF(@name::text, ''), i.name) IS NULL OR i.name LIKE '%' || COALESCE(NULLIF(@name::text, ''), i.name) || '%')
--   AND (COALESCE(@score::int, i.score) IS NULL OR i.score >= COALESCE(@score::int, i.score))
--   AND (COALESCE(NULLIF(@type, '')::"ItemType", i.type) IS NULL OR i.type = COALESCE(NULLIF(@type, '')::"ItemType", i.type))
--   AND (COALESCE(NULLIF(@max_price::int, 0), CASE WHEN i.discount_active = true THEN i.discount_price ELSE i.price END) IS NULL OR
--        CASE
--          WHEN i.discount_active = true THEN i.discount_price
--          ELSE i.price
--        END <= COALESCE(NULLIF(@max_price::int, 0), CASE
--                              WHEN i.discount_active = true THEN i.discount_price
--                              ELSE i.price
--                            END))
-- ORDER BY i.score DESC, s.score DESC;
--
-- -- name: IsItemOwner :one
-- SELECT EXISTS(
--     SELECT 1
--     FROM item i
--              INNER JOIN store s ON s.id = i.store_id
--     WHERE i.id = $1
--       AND s.owner_id = $2
-- );
--
-- -- name: SetItemImage :exec
-- UPDATE item
-- SET
--     image = $2
-- WHERE id = $1;