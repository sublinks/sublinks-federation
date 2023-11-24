-- +goose Up
CREATE TABLE users (
    id  INTEGER PRIMARY KEY,
    name CHAR(50) NOT NULL,
    public_key TEXT NOT NULL,
    private_key TEXT NOT NULL,
    unique (name)
);

-- +goose Down
DROP TABLE users;