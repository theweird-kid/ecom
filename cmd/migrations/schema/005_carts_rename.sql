-- +goose Up
ALTER TABLE carts RENAME COLUMN "total_price" TO price;

-- +goose Down
ALTER TABLE carts RENAME COLUMN price TO "total_price";