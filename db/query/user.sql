-- name: CreateUser :one
INSERT INTO users (
  code,
  username,
  first_name,
  last_name,
  email,
  phone,
  description,
  hashed_password
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: ListUsers :many
SELECT id,
code,
first_name,
last_name,
email,
phone,
description,
created_at,
updated_at
FROM users
WHERE deleted IS False
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetUser :one
SELECT id,
code,
first_name,
last_name,
username,
email,
phone,
description
FROM users
WHERE id = $1 
AND deleted IS False
LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1
AND deleted IS False
LIMIT 1;

-- name: CountUser :one
SELECT COUNT(1) AS "UserCount" 
FROM users 
WHERE deleted IS False
LIMIT 1;

-- name: CountUserForGenerateCode :one
SELECT COUNT(1) AS "UserCount" 
FROM users
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
  first_name = COALESCE(sqlc.narg(first_name), first_name),
  last_name = COALESCE(sqlc.narg(last_name), last_name),
  email = COALESCE(sqlc.narg(email), email),
  phone = COALESCE(sqlc.narg(phone), phone),
  description = COALESCE(sqlc.narg(description), description),
  -- hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  -- password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
  updated_at = NOW()::TIMESTAMP
WHERE
  id = sqlc.arg(user_id) AND deleted IS False
RETURNING *;

-- name: DeleteUser :exec
UPDATE users
SET deleted = True,
updated_at = NOW()
WHERE id = $1
AND deleted IS False;

-- name: GetUserHashedPassword :one
SELECT hashed_password
FROM users
WHERE id = $1
AND deleted IS False;

-- name: UpdateUserPassword :exec
UPDATE users
SET
  hashed_password = $2,
  password_changed_at = NOW()::TIMESTAMP,
  updated_at = NOW()::TIMESTAMP
WHERE
  id = $1
  AND deleted IS False;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
AND deleted IS False;
