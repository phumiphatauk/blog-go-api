-- name: CreateSession :one
INSERT INTO sessions (
  id,
  user_id,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetSession :one
SELECT s.*, u.username
FROM sessions AS s
INNER JOIN users AS u
ON s.user_id = u.id
WHERE s.id = $1 LIMIT 1;
