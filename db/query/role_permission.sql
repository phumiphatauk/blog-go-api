-- name: GetRolePermissionByRoleIdAndPermissionId :one
SELECT
id, role_id, permission_id, created_at, updated_at, deleted
FROM "role_permission"
WHERE deleted IS False 
AND role_id = $1
AND permission_id = $2
LIMIT 1;

-- name: CreateRolePermission :one
INSERT INTO "role_permission" (role_id, permission_id, created_at, updated_at, deleted)
VALUES ($1, $2, NOW()::TIMESTAMP, NULL, FALSE)
RETURNING id, role_id, permission_id, created_at, updated_at, deleted;

-- name: DeleteRolePermission :exec
UPDATE "role_permission"
SET updated_at = NOW()::TIMESTAMP,
deleted = TRUE
WHERE deleted IS False
AND role_id = $1
AND permission_id = $2;

-- name: DeleteRolePermissionByRoleId :exec
UPDATE "role_permission"
SET updated_at = NOW()::TIMESTAMP,
deleted = TRUE
WHERE role_id = $1;
