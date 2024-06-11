-- +goose Up
CREATE TABLE orders(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    prod_id UUID NOT NULL REFERENCES products(id),
    created_at TIMESTAMP NOT NULL,
    price INTEGER NOT NULL,
    quantity INTEGER NOT NULL
);

-- +goose Down
DROP TABLE orders;