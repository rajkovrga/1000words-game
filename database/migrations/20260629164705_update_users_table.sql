-- +goose Up

PRAGMA foreign_keys = OFF;

CREATE TABLE users_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

INSERT INTO users_new (
    id,
    email,
    password,
    created_at,
    updated_at,
    deleted_at
)
SELECT
    id,
    username AS email,
    '' AS password,
    created_at,
    CURRENT_TIMESTAMP AS updated_at,
    NULL AS deleted_at
FROM users;

DROP TABLE users;

ALTER TABLE users_new RENAME TO users;

CREATE INDEX IF NOT EXISTS idx_users_deleted_at
ON users(deleted_at);

PRAGMA foreign_keys = ON;

-- +goose Down

PRAGMA foreign_keys = OFF;

DROP INDEX IF EXISTS idx_users_deleted_at;

CREATE TABLE users_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users_old (
    id,
    username,
    created_at
)
SELECT
    id,
    email AS username,
    created_at
FROM users;

DROP TABLE users;

ALTER TABLE users_old RENAME TO users;

PRAGMA foreign_keys = ON;
