-- name: CreateResetPassword :one
INSERT INTO reset_password (
    user_id,
    token,
    expires_at
) VALUES (
$1, $2, $3
) RETURNING *;

-- name: GetResetPasswordByToken :one
SELECT * FROM reset_password
WHERE token = $1
AND used IS FALSE;

-- name: UseResetPassword :exec
UPDATE reset_password
SET used = TRUE
WHERE token = $1;
