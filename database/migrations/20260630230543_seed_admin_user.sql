-- +goose Up

INSERT INTO users (
    email,
    password,
    created_at,
    updated_at,
    deleted_at
)
SELECT
    'admin@vrga.dev',
    '$2a$10$wQ4T6lP6C2LDhC./uCUrxegc7SYcmmXTTGNY.ScpTEWmbkvMTNhRm',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    NULL
WHERE NOT EXISTS (
    SELECT 1
    FROM users
    WHERE email = 'admin@vrga.dev'
);

INSERT OR IGNORE INTO user_roles (
    user_id,
    role_id
)
SELECT
    users.id,
    roles.id
FROM users
JOIN roles
    ON roles.code = 'admin'
WHERE users.email = 'admin@vrga.dev'
  AND users.deleted_at IS NULL
  AND roles.deleted_at IS NULL;


-- +goose Down

DELETE FROM user_roles
WHERE user_id IN (
    SELECT id
    FROM users
    WHERE email = 'admin@vrga.dev'
);

DELETE FROM users
WHERE email = 'admin@vrga.dev';