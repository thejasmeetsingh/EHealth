-- +goose Up

CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    modified_at TIMESTAMP NOT NULL,
    email VARCHAR(50) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    name VARCHAR(50),
    is_end_user BOOLEAN NOT NULL  -- This will depict wheather the user is end user or a medical facility root user.
);

-- +goose Down
DROP TABLE users;