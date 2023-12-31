// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: users.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, modified_at, email, password, is_end_user) 
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at, modified_at, email, password, name, is_end_user
`

type CreateUserParams struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	ModifiedAt time.Time
	Email      string
	Password   string
	IsEndUser  bool
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.CreatedAt,
		arg.ModifiedAt,
		arg.Email,
		arg.Password,
		arg.IsEndUser,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.ModifiedAt,
		&i.Email,
		&i.Password,
		&i.Name,
		&i.IsEndUser,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id=$1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, modified_at, email, password, name, is_end_user FROM users WHERE email=$1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.ModifiedAt,
		&i.Email,
		&i.Password,
		&i.Name,
		&i.IsEndUser,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, created_at, modified_at, email, password, name, is_end_user FROM users WHERE id=$1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.ModifiedAt,
		&i.Email,
		&i.Password,
		&i.Name,
		&i.IsEndUser,
	)
	return i, err
}

const updateUserDetails = `-- name: UpdateUserDetails :one
UPDATE users SET email=$1, name=$2 
WHERE id=$3
RETURNING id, created_at, modified_at, email, password, name, is_end_user
`

type UpdateUserDetailsParams struct {
	Email string
	Name  sql.NullString
	ID    uuid.UUID
}

func (q *Queries) UpdateUserDetails(ctx context.Context, arg UpdateUserDetailsParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserDetails, arg.Email, arg.Name, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.ModifiedAt,
		&i.Email,
		&i.Password,
		&i.Name,
		&i.IsEndUser,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :one
UPDATE users SET password=$1
WHERE id=$2
RETURNING id, created_at, modified_at, email, password, name, is_end_user
`

type UpdateUserPasswordParams struct {
	Password string
	ID       uuid.UUID
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserPassword, arg.Password, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.ModifiedAt,
		&i.Email,
		&i.Password,
		&i.Name,
		&i.IsEndUser,
	)
	return i, err
}
