CREATE TABLE roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE role_permission (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    FOREIGN KEY (role_id) REFERENCES roles(id),
    FOREIGN KEY (permission_id) REFERENCES permissions(id)
);

CREATE INDEX idx_role_permission_role_id ON role_permission(role_id);
CREATE INDEX idx_role_permission_permission_id ON role_permission(permission_id);

CREATE TABLE languages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL UNIQUE,
    language_name TEXT NOT NULL
);

CREATE TABLE words (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    number INTEGER NOT NULL
);

CREATE INDEX idx_words_number ON words(number);

CREATE TABLE word_translations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    word_id INTEGER NOT NULL,
    language_id INTEGER NOT NULL,
    word_value TEXT NOT NULL,
    FOREIGN KEY (word_id) REFERENCES words(id),
    FOREIGN KEY (language_id) REFERENCES languages(id),
    UNIQUE(word_id, language_id)
);

CREATE INDEX idx_word_translations_word_id ON word_translations(word_id);
CREATE INDEX idx_word_translations_language_id ON word_translations(language_id);

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role_id INTEGER NOT NULL,
    native_language_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (role_id) REFERENCES roles(id),
    FOREIGN KEY (native_language_id) REFERENCES languages(id)
);

CREATE INDEX idx_users_role_id ON users(role_id);
CREATE INDEX idx_users_native_language_id ON users(native_language_id);

CREATE TABLE level_results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    target_language_id INTEGER NOT NULL,
    level INTEGER NOT NULL,
    success INTEGER DEFAULT 0,
    level_completed_at DATETIME,
    total_completed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (target_language_id) REFERENCES languages(id),
    UNIQUE (user_id, target_language_id, level)
);

CREATE INDEX idx_level_results_user_language ON level_results(user_id, target_language_id);
CREATE INDEX idx_level_results_level ON level_results(level);
