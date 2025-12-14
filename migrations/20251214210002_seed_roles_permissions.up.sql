-- ROLES
INSERT INTO roles (id, name) VALUES
(1, 'admin'),
(2, 'user');

-- PERMISSIONS
INSERT INTO permissions (name) VALUES
-- user
('user.view.self'),
('user.update.self'),

-- languages
('languages.view'),
('languages.create'),
('languages.update'),
('languages.delete'),

-- words
('words.view'),
('words.create'),
('words.update'),
('words.delete'),

-- translations
('translations.view'),
('translations.create'),
('translations.update'),
('translations.delete'),

-- practice
('practice.start'),
('practice.submit'),
('practice.view.stats'),

-- admin
('users.view'),
('users.create'),
('users.update'),
('users.delete'),
('roles.manage'),
('permissions.manage');

INSERT INTO role_permission (role_id, permission_id)
SELECT 1, id FROM permissions;

INSERT INTO role_permission (role_id, permission_id)
SELECT 2, id FROM permissions
WHERE name IN (
    'user.view.self',
    'user.update.self',
    'languages.view',
    'words.view',
    'translations.view',
    'practice.start',
    'practice.submit',
    'practice.view.stats'
);
