DROP INDEX IF EXISTS idx_level_results_level;
DROP INDEX IF EXISTS idx_level_results_user_language;

DROP INDEX IF EXISTS idx_users_native_language_id;
DROP INDEX IF EXISTS idx_users_role_id;

DROP INDEX IF EXISTS idx_word_translations_language_id;
DROP INDEX IF EXISTS idx_word_translations_word_id;

DROP INDEX IF EXISTS idx_words_number;

DROP INDEX IF EXISTS idx_role_permission_permission_id;
DROP INDEX IF EXISTS idx_role_permission_role_id;

DROP TABLE IF EXISTS level_results;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS word_translations;
DROP TABLE IF EXISTS words;
DROP TABLE IF EXISTS languages;
DROP TABLE IF EXISTS role_permission;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
