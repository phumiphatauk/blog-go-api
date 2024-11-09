-- name: GetAllBlog :many
SELECT
id, title, content, image, url, created_at, updated_at
FROM blog
WHERE deleted IS FALSE
AND LOWER(title) LIKE '%' || LOWER($1) || '%'
ORDER BY created_at DESC
OFFSET $2
LIMIT $3;

-- name: CountAllBlog :one
SELECT COUNT(1) AS count
FROM blog
WHERE deleted IS FALSE
AND LOWER(title) LIKE '%' || LOWER($1) || '%'
LIMIT 1;

-- name: GetAllBlogWithTag :many
SELECT DISTINCT
b.id, b.title, b.content, b.image, b.url, b.created_at, b.updated_at
FROM blog b
LEFT JOIN blog_tag bt ON b.id = bt.blog_id AND bt.deleted IS FALSE
LEFT JOIN tag t ON bt.tag_id = t.id AND t.deleted IS FALSE
WHERE b.deleted IS FALSE
AND LOWER(b.title) LIKE '%' || LOWER($1) || '%'
AND LOWER(t.name) = LOWER($2)
ORDER BY b.created_at DESC
OFFSET $3
LIMIT $4;

-- name: CountAllBlogWithTag :one
SELECT COUNT(DISTINCT b.id) AS count
FROM blog b
LEFT JOIN blog_tag bt ON b.id = bt.blog_id AND bt.deleted IS FALSE
LEFT JOIN tag t ON bt.tag_id = t.id AND t.deleted IS FALSE
WHERE b.deleted IS FALSE
AND LOWER(b.title) LIKE '%' || LOWER($1) || '%'
AND LOWER(t.name) = LOWER($2)
LIMIT 1;

-- name: GetBlogById :one
SELECT
id, title, content, image, url, created_at, updated_at
FROM blog
WHERE deleted IS FALSE
AND id = $1
LIMIT 1;

-- name: GetBlogByUrl :one
SELECT
id, title, content, image, url, created_at, updated_at
FROM blog
WHERE deleted IS FALSE
AND url = $1
LIMIT 1;

-- name: CreateBlog :one
INSERT INTO blog
(title, content, image, url, created_at)
VALUES ($1, $2, $3, $4, NOW()::TIMESTAMPTZ)
RETURNING *;

-- name: UpdateBlog :exec
UPDATE blog
SET title = $2,
content = $3,
image = $4,
url = $5,
updated_at = NOW()::TIMESTAMPTZ
WHERE deleted IS FALSE
AND id = $1;

-- name: DeleteBlog :exec
UPDATE blog
SET updated_at = NOW()::TIMESTAMPTZ,
deleted = TRUE
WHERE deleted IS FALSE
AND id = $1;
