INSERT INTO role_permission (role_id, permission_id, created_at, updated_at, deleted) 
SELECT 1, id, NOW()::TIMESTAMP, NULL, FALSE
FROM permission;
