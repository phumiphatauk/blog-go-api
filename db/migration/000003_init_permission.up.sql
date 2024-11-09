-- Dashboards
INSERT INTO permission_group (id, name) VALUES (1, 'Dashboards');
INSERT INTO permission (id, code, name, permission_group_id) VALUES (1, 'view_dashboard_analytics', 'View Dashboard Analytics', 1);

-- Users
INSERT INTO permission_group (id, name) VALUES (2, 'Users');
INSERT INTO permission (id, code, name, permission_group_id) VALUES (2, 'view_user', 'View User', 2);
INSERT INTO permission (id, code, name, permission_group_id) VALUES (3, 'edit_user', 'Edit User', 2);

-- Roles
INSERT INTO permission_group (id, name) VALUES (3, 'Roles');
INSERT INTO permission (id, code, name, permission_group_id) VALUES (4, 'view_role', 'View Role', 3);
INSERT INTO permission (id, code, name, permission_group_id) VALUES (5, 'edit_role', 'Edit Role', 3);
