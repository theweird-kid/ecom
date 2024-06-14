-- +goose Up
CREATE TABLE carts(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    prod_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    updated_at TIMESTAMP NOT NULL,
    total_price INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    UNIQUE(user_id, prod_id)
);

-- +goose Down
DROP TABLE carts;