-- name: CreateUserRole :exec
INSERT INTO user_role 
(user_id, role_id, created_at) 
VALUES 
($1, $2, NOW()::TIMESTAMP);

-- name: DeleteUserRoleByUserId :exec
UPDATE user_role
SET deleted = True,
updated_at = NOW()::TIMESTAMP
WHERE user_id = $1;
