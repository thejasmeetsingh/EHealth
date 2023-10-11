-- name: CreateUser :one
INSERT INTO users (id, created_at, modified_at, email, password, is_end_user) 
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id=$1;