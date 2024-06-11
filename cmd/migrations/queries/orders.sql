-- name: GetOrders :many
SELECT * FROM orders
WHERE user_id = $1;

-- name: OrderItem :one
INSERT INTO orders(id, user_id, prod_id, created_at, price, quantity)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

