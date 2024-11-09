-- Blogs
INSERT INTO permission_group (id, name) VALUES (4, 'Blogs');
INSERT INTO permission (id, code, name, permission_group_id) VALUES (6, 'view_blog', 'View Blog', 4);
INSERT INTO permission (id, code, name, permission_group_id) VALUES (7, 'edit_blog', 'Edit Blog', 4);

-- Tags
INSERT INTO permission_group (id, name) VALUES (5, 'Tags');
INSERT INTO permission (id, code, name, permission_group_id) VALUES (8, 'view_tag', 'View Tag', 5);
INSERT INTO permission (id, code, name, permission_group_id) VALUES (9, 'edit_tag', 'Edit Tag', 5);
