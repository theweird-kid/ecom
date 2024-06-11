-- name: CreateProduct :one
INSERT INTO products(id, created_at, updated_at, name, description, image, price, quantity)
VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetProducts :many
SELECT * FROM products;

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1;

-- name: UpdateProduct :one
UPDATE products
SET updated_at = $1,
    name = $2,
    description = $3,
    image = $4,
    price = $5,
    quantity = $6
WHERE id = $7
RETURNING *;

-- name: DeleteProduct :one
DELETE FROM products
WHERE id = $1
RETURNING *;