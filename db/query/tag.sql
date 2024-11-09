-- name: GetAllTag :many
SELECT
id, name
FROM tag
WHERE deleted IS FALSE
AND LOWER(name) LIKE '%' || LOWER($1) || '%'
ORDER BY name ASC
OFFSET $2
LIMIT $3;

-- name: CountAllTag :one
SELECT COUNT(1) AS count
FROM tag
WHERE deleted IS FALSE
AND LOWER(name) LIKE '%' || LOWER($1) || '%' 
LIMIT 1;

-- name: GetTagById :one
SELECT
id, name
FROM tag
WHERE deleted IS FALSE
AND id = $1
LIMIT 1;

-- name: CreateTag :one
INSERT INTO tag
(name, created_at)
VALUES ($1, NOW()::TIMESTAMPTZ)
RETURNING *;

-- name: UpdateTag :exec
UPDATE tag
SET name = $2,
updated_at = NOW()::TIMESTAMPTZ
WHERE deleted IS FALSE
AND id = $1;

-- name: DeleteTag :exec
UPDATE tag
SET updated_at = NOW()::TIMESTAMPTZ,
deleted = TRUE
WHERE deleted IS FALSE
AND id = $1;

-- name: GetTagByCreatedAt :many
SELECT
id, name
FROM tag
WHERE deleted IS FALSE
AND created_at = $1
ORDER BY name ASC;
