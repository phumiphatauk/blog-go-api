-- name: GetPermissionByPermissionGroupId :many
SELECT
id, name
FROM permission
WHERE permission_group_id = $1;

-- name: GetPermissionByPermissionGroupIdAndRoleId :many
SELECT
p.id,
p.name,
CASE WHEN rp.id IS NOT NULL THEN TRUE ELSE FALSE END AS is_assigned
FROM permission p
LEFT JOIN role_permission rp ON p.id = rp.permission_id AND rp.deleted IS FALSE AND rp.role_id = $1
WHERE permission_group_id = $2;

-- name: GetPermissionByUserId :many
SELECT 
DISTINCT
p.code
FROM permission p
INNER JOIN permission_group pg ON p.permission_group_id = pg.id
INNER JOIN role_permission rp ON p.id = rp.permission_id AND rp.deleted IS FALSE
INNER JOIN role r ON rp.role_id = r.id AND r.deleted IS FALSE
INNER JOIN user_role ur ON r.id = ur.role_id AND ur.deleted IS FALSE
INNER JOIN users u ON ur.user_id = u.id AND u.deleted IS FALSE
WHERE u.id = $1;
