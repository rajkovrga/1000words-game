-- +goose Up

INSERT OR IGNORE INTO roles (
    code,
    name,
    description
)
VALUES
    ('user', 'User', 'Regular application user'),
    ('admin', 'Admin', 'Full application administrator'),
    ('moderator', 'Moderator', 'Dashboard user focused on managing learning content');


INSERT OR IGNORE INTO permissions (
    code,
    name,
    description
)
VALUES
    ('me.read', 'Read own profile', 'User can read own profile'),

    ('languages.read', 'Read languages', 'User can read available languages'),
    ('levels.read', 'Read levels', 'User can read available levels'),

    ('practice.start', 'Start practice', 'User can start practice mode'),

    ('progress.read', 'Read own progress', 'User can read own learning progress'),
    ('progress.create', 'Create own progress', 'User can create new learning progress'),
    ('progress.update', 'Update own progress', 'User can update own learning progress'),

    ('game.start', 'Start game', 'User can start game mode'),
    ('game.finish', 'Finish game', 'User can finish game attempt'),
    ('attempts.read', 'Read own attempts', 'User can read own game attempts'),

    ('dashboard.access', 'Access dashboard', 'User can access dashboard'),

    ('dashboard.stats.read', 'Read dashboard stats', 'User can read dashboard statistics'),

    ('users.read', 'Read users', 'User can read users in dashboard'),
    ('users.create', 'Create users', 'User can create users in dashboard'),
    ('users.update', 'Update users', 'User can update users in dashboard'),
    ('users.delete', 'Delete users', 'User can delete users in dashboard'),

    ('roles.read', 'Read roles', 'User can read roles'),
    ('roles.create', 'Create roles', 'User can create roles'),
    ('roles.update', 'Update roles', 'User can update roles'),
    ('roles.delete', 'Delete roles', 'User can delete roles'),

    ('permissions.read', 'Read permissions', 'User can read permissions'),
    ('permissions.create', 'Create permissions', 'User can create permissions'),
    ('permissions.update', 'Update permissions', 'User can update permissions'),
    ('permissions.delete', 'Delete permissions', 'User can delete permissions'),

    ('words.read', 'Read words', 'User can read words'),
    ('words.create', 'Create words', 'User can create new words'),
    ('words.update', 'Update words', 'User can update existing words'),
    ('words.delete', 'Delete words', 'User can delete words'),
    ('words.import', 'Import words', 'User can import words in bulk'),

    ('word_translations.read', 'Read word translations', 'User can read word translations'),
    ('word_translations.create', 'Create word translations', 'User can create word translations'),
    ('word_translations.update', 'Update word translations', 'User can update word translations'),
    ('word_translations.delete', 'Delete word translations', 'User can delete word translations'),

    ('admin.settings.read', 'Read admin settings', 'User can read admin settings'),
    ('admin.settings.update', 'Update admin settings', 'User can update admin settings'),

    ('auth_tokens.read', 'Read auth tokens', 'User can read authentication tokens'),
    ('auth_tokens.delete', 'Delete auth tokens', 'User can delete authentication tokens');


WITH allowed_permissions(role_code, permission_code) AS (
    VALUES
        ('user', 'me.read'),
        ('user', 'languages.read'),
        ('user', 'levels.read'),
        ('user', 'practice.start'),
        ('user', 'progress.read'),
        ('user', 'progress.create'),
        ('user', 'progress.update'),
        ('user', 'game.start'),
        ('user', 'game.finish'),
        ('user', 'attempts.read')
)
INSERT OR IGNORE INTO role_permissions (
    role_id,
    permission_id
)
SELECT
    roles.id,
    permissions.id
FROM allowed_permissions
JOIN roles
    ON roles.code = allowed_permissions.role_code
JOIN permissions
    ON permissions.code = allowed_permissions.permission_code;


WITH allowed_permissions(role_code, permission_code) AS (
    VALUES
        ('moderator', 'dashboard.access'),
        ('moderator', 'dashboard.stats.read'),

        ('moderator', 'languages.read'),
        ('moderator', 'levels.read'),

        ('moderator', 'words.read'),
        ('moderator', 'words.create'),
        ('moderator', 'words.update'),
        ('moderator', 'words.delete'),
        ('moderator', 'words.import'),

        ('moderator', 'word_translations.read'),
        ('moderator', 'word_translations.create'),
        ('moderator', 'word_translations.update'),
        ('moderator', 'word_translations.delete')
)
INSERT OR IGNORE INTO role_permissions (
    role_id,
    permission_id
)
SELECT
    roles.id,
    permissions.id
FROM allowed_permissions
JOIN roles
    ON roles.code = allowed_permissions.role_code
JOIN permissions
    ON permissions.code = allowed_permissions.permission_code;


INSERT OR IGNORE INTO role_permissions (
    role_id,
    permission_id
)
SELECT
    roles.id,
    permissions.id
FROM roles
JOIN permissions
WHERE roles.code = 'admin';


-- +goose Down

DELETE FROM role_permissions
WHERE role_id IN (
    SELECT id
    FROM roles
    WHERE code IN ('user', 'admin', 'moderator')
);

DELETE FROM user_roles
WHERE role_id IN (
    SELECT id
    FROM roles
    WHERE code IN ('user', 'admin', 'moderator')
);

DELETE FROM permissions
WHERE code IN (
    'me.read',

    'languages.read',
    'levels.read',

    'practice.start',

    'progress.read',
    'progress.create',
    'progress.update',

    'game.start',
    'game.finish',
    'attempts.read',

    'dashboard.access',
    'dashboard.stats.read',

    'users.read',
    'users.create',
    'users.update',
    'users.delete',

    'roles.read',
    'roles.create',
    'roles.update',
    'roles.delete',

    'permissions.read',
    'permissions.create',
    'permissions.update',
    'permissions.delete',

    'words.read',
    'words.create',
    'words.update',
    'words.delete',
    'words.import',

    'word_translations.read',
    'word_translations.create',
    'word_translations.update',
    'word_translations.delete',

    'admin.settings.read',
    'admin.settings.update',

    'auth_tokens.read',
    'auth_tokens.delete'
);

DELETE FROM roles
WHERE code IN ('user', 'admin', 'moderator');