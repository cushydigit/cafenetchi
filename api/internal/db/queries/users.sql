-- name: ListUsers :many
SELECT * 
FROM users;

-- name: GetUserByID :one
SELECT * 
FROM users
WHERE id = $1;

-- name: GetUserByPhone :one
SELECT *
FROM users 
WHERE phone = $1;

-- name: CreateUser :one
INSERT INTO users (phone, full_name, is_verified, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUser :one
UPDATE users 
SET full_name = $2, 
    avatar_url = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;