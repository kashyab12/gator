-- +goose Up
CREATE TABLE users (
    uuid INTEGER PRIMARY KEY,
    created_at TIMESTAM NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL
)

-- +goose Down
DROP TABLE users;