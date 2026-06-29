-- +goose Up
INSERT OR IGNORE INTO levels 
(id, level_number, name, words_required, max_wrong_answers) 
VALUES
(1, 1, 'Level 1', 100, 3),
(2, 2, 'Level 2', 100, 3),
(3, 3, 'Level 3', 100, 3),
(4, 4, 'Level 4', 100, 3),
(5, 5, 'Level 5', 100, 3),
(6, 6, 'Level 6', 100, 3),
(7, 7, 'Level 7', 100, 3),
(8, 8, 'Level 8', 100, 3),
(9, 9, 'Level 9', 100, 3),
(10, 10, 'Level 10', 100, 3);

delete from languages;

INSERT OR IGNORE INTO languages 
(id, name, code) 
VALUES
(1, 'English', 'en'),
(2, 'Serbian', 'sr'),
(3, 'Spanish', 'es');

-- +goose Down
delete from levels;
delete from languages;

