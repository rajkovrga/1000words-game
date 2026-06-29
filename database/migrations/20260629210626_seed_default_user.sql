-- +goose Up

INSERT INTO users (
    email,
    password,
    created_at,
    updated_at,
    deleted_at
)
SELECT
    'rajko@vrga.dev',
    '$2y$10$q3.hfqDJwf0xrahvX/1le.nBxiAB36g.zS6DnvlfSWFZqyC3ZIwH.',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    NULL
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'rajko@vrga.dev'
);

-- +goose Down

DELETE FROM users
WHERE email = 'rajko@vrga.dev';