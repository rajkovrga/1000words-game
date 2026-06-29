-- +goose Up
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS levels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    level_number INTEGER NOT NULL UNIQUE,
    name TEXT NOT NULL,
    words_required INTEGER NOT NULL DEFAULT 100,
    max_wrong_answers INTEGER NOT NULL DEFAULT 3,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS languages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    code TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS words (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    level_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (level_id) REFERENCES levels(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS word_languages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    word_id INTEGER NOT NULL,
    language_id INTEGER NOT NULL,
    text TEXT NOT NULL,
    normalized_text TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (word_id) REFERENCES words(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    FOREIGN KEY (language_id) REFERENCES languages(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    UNIQUE(word_id, language_id)
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_language_progress (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    target_language_id INTEGER NOT NULL,
    native_language_id INTEGER NOT NULL,
    current_level_id INTEGER NOT NULL,
    total_attempts INTEGER NOT NULL DEFAULT 0,
    total_passed INTEGER NOT NULL DEFAULT 0,
    total_failed INTEGER NOT NULL DEFAULT 0,
    total_correct_answers INTEGER NOT NULL DEFAULT 0,
    total_wrong_answers INTEGER NOT NULL DEFAULT 0,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    FOREIGN KEY (target_language_id) REFERENCES languages(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    FOREIGN KEY (native_language_id) REFERENCES languages(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    FOREIGN KEY (current_level_id) REFERENCES levels(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    UNIQUE(user_id, target_language_id, native_language_id),

    CHECK(target_language_id != native_language_id)
);

CREATE TABLE IF NOT EXISTS user_attempts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    target_language_id INTEGER NOT NULL,
    native_language_id INTEGER NOT NULL,
    level_id INTEGER NOT NULL,
    correct_count INTEGER NOT NULL DEFAULT 0,
    wrong_count INTEGER NOT NULL DEFAULT 0,
    total_questions INTEGER NOT NULL DEFAULT 0,
    passed INTEGER NOT NULL DEFAULT 0,
    started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finished_at DATETIME,

    FOREIGN KEY (user_id) REFERENCES users(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    FOREIGN KEY (target_language_id) REFERENCES languages(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    FOREIGN KEY (native_language_id) REFERENCES languages(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    FOREIGN KEY (level_id) REFERENCES levels(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CHECK(target_language_id != native_language_id),
    CHECK(passed IN (0, 1))
);

CREATE INDEX IF NOT EXISTS idx_words_level_id 
ON words(level_id);

CREATE INDEX IF NOT EXISTS idx_word_languages_word_id 
ON word_languages(word_id);

CREATE INDEX IF NOT EXISTS idx_word_languages_language_id 
ON word_languages(language_id);

CREATE INDEX IF NOT EXISTS idx_word_languages_language_text 
ON word_languages(language_id, normalized_text);

CREATE INDEX IF NOT EXISTS idx_user_language_progress_user_id 
ON user_language_progress(user_id);

CREATE INDEX IF NOT EXISTS idx_user_language_progress_languages 
ON user_language_progress(target_language_id, native_language_id);

CREATE INDEX IF NOT EXISTS idx_user_attempts_user_id 
ON user_attempts(user_id);

CREATE INDEX IF NOT EXISTS idx_user_attempts_user_level 
ON user_attempts(user_id, level_id);

CREATE INDEX IF NOT EXISTS idx_user_attempts_languages_level 
ON user_attempts(target_language_id, native_language_id, level_id);

-- +goose Down

DROP INDEX IF EXISTS idx_user_attempts_languages_level;
DROP INDEX IF EXISTS idx_user_attempts_user_level;
DROP INDEX IF EXISTS idx_user_attempts_user_id;

DROP INDEX IF EXISTS idx_user_language_progress_languages;
DROP INDEX IF EXISTS idx_user_language_progress_user_id;

DROP INDEX IF EXISTS idx_word_languages_language_text;
DROP INDEX IF EXISTS idx_word_languages_language_id;
DROP INDEX IF EXISTS idx_word_languages_word_id;

DROP INDEX IF EXISTS idx_words_level_id;

DROP TABLE IF EXISTS user_attempts;
DROP TABLE IF EXISTS user_language_progress;
DROP TABLE IF EXISTS word_languages;
DROP TABLE IF EXISTS words;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS languages;
DROP TABLE IF EXISTS levels;
