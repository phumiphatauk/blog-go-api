-- name: GetBlogTagByBlogId :many
SELECT
bt.id,
bt.blog_id,
bt.tag_id,
t.name,
bt.deleted
FROM blog_tag bt
INNER JOIN tag t ON bt.tag_id = t.id AND t.deleted IS FALSE
WHERE bt.deleted IS FALSE
AND blog_id = $1;

-- name: GetBlogTagByBlogIdAndTagId :one
SELECT
id, blog_id, tag_id, created_at, updated_at, deleted
FROM blog_tag
WHERE deleted IS FALSE
AND blog_id = $1
AND tag_id = $2
LIMIT 1;

-- name: CreateBlogTag :exec
INSERT INTO blog_tag
(blog_id, tag_id, created_at)
VALUES
($1, $2, NOW()::TIMESTAMP);

-- name: DeleteBlogTag :exec
UPDATE blog_tag
SET deleted = True,
updated_at = NOW()::TIMESTAMP
WHERE blog_id = $1
AND tag_id = $2;

-- name: DeleteBlogTagByBlogId :exec
UPDATE blog_tag
SET deleted = True,
updated_at = NOW()::TIMESTAMP
WHERE blog_id = $1;
