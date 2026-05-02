-- name: CreateUser :one
INSERT INTO users(
email,password_hash,name) VALUES ($1,$2,$3) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;


-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;


-- name: AddUserBalance :one
UPDATE users
SET balance_cents = balance_cents + $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;


-- name: DeductUserBalance :one
UPDATE users
SET balance_cents = balance_cents - $2,
    updated_at = NOW()
WHERE id = $1
  AND balance_cents >= $2
RETURNING *;



-- name: ListUsers :many
SELECT *
FROM users
ORDER BY id DESC
LIMIT $1 OFFSET $2;
