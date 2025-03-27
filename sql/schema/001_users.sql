-- +goose Up
CREATE TABLE users (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    email VARCHAR(255) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users; 