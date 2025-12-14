INSERT INTO languages (id, code, language_name) VALUES
(1, 'sr', 'Serbian'),
(2, 'en', 'English');

INSERT INTO words (id, number) VALUES
(1, 1),
(2, 2),
(3, 3),
(4, 4),
(5, 5);

INSERT INTO word_translations (word_id, language_id, word_value) VALUES
(1, 1, 'mačka'),
(1, 2, 'cat'),
(2, 1, 'pas'),
(2, 2, 'dog'),
(3, 1, 'kuća'),
(3, 2, 'house'),
(4, 1, 'voda'),
(4, 2, 'water'),
(5, 1, 'hleb'),
(5, 2, 'bread');

-- password: test123
INSERT INTO users (
    id,
    email,
    password_hash,
    role_id,
    native_language_id,
    created_at
) VALUES (
    1,
    'rajko@vrga.dev',
    '$2a$10$7QZ6pZ7JqD1N7X9zq9r5TOV0kz8Dkq0U8T7dP4E5YxCzY2z0a7C3K',
    2,
    1,
    datetime('now')
);
