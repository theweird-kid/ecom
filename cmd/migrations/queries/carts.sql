-- name: GetCart :many
SELECT * FROM carts
WHERE user_id = $1;

-- name: AddToCart :one
WITH valid_quantity AS (
    SELECT id FROM products
    WHERE id = $3 AND quantity >= $6
)
INSERT INTO carts(id, user_id, prod_id, updated_at, total_price, quantity) 
SELECT $1, $2, $3, $4, $5, $6
FROM valid_quantity
WHERE EXISTS (
    SELECT 1 FROM valid_quantity
)
RETURNING *;

-- name: UpdateCart :one
UPDATE carts
SET quantity = $1, updated_at = $2,
total_price = $3
WHERE id = $4
RETURNING *;

-- name: DeleteCartItem :exec
DELETE FROM carts
WHERE id = $1;