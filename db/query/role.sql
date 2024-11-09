-- name: GetAllRole :many
SELECT 
id, name
FROM "role"
WHERE deleted IS False
AND LOWER(name) LIKE '%' || LOWER($1) || '%'
ORDER BY name ASC
OFFSET $2
LIMIT $3;

-- name: CountAllRole :one
SELECT COUNT(1) AS count 
FROM "role"
WHERE deleted IS False
AND LOWER(name) LIKE '%' || LOWER($1) || '%'
LIMIT 1;

-- name: GetRoleById :one
SELECT
id, name
FROM "role"
WHERE deleted IS False 
AND id = $1
LIMIT 1;

-- name: CreateRole :one
INSERT INTO "role" (name, created_at)
VALUES ($1, NOW()::TIMESTAMP)
RETURNING *;

-- name: UpdateRole :exec
UPDATE "role"
SET name = $2,
updated_at = NOW()::TIMESTAMP
WHERE deleted IS False 
AND id = $1;

-- name: DeleteRole :exec
UPDATE "role"
SET updated_at = NOW()::TIMESTAMP,
deleted = TRUE
WHERE deleted IS False 
AND id = $1;

-- name: GetRoleByUserId :many
SELECT r.id
, r.name
FROM "role" r
INNER JOIN "user_role" ur ON r.id = ur.role_id AND ur.deleted IS false
WHERE r.deleted IS false
AND ur.user_id = $1
ORDER BY r.name ASC;

-- name: GetRoleForDropDownList :many
SELECT 
id, name
FROM "role"
WHERE deleted IS False
ORDER BY name ASC;