DELETE FROM users WHERE id = 1;

DELETE FROM word_translations
WHERE word_id IN (1, 2, 3, 4, 5);

DELETE FROM words
WHERE id IN (1, 2, 3, 4, 5);

DELETE FROM languages
WHERE id IN (1, 2);
