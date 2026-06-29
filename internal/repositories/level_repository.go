package repositories

import (
	"database/sql"

	dbModels "1000words-game/models/db"
)

type LevelRepository struct {
	db *sql.DB
}

func NewLevelRepository(db *sql.DB) *LevelRepository {
	return &LevelRepository{
		db: db,
	}
}

func (r *LevelRepository) GetAll() ([]dbModels.Level, error) {
	query := `
		SELECT
			id,
			level_number,
			name,
			words_required,
			max_wrong_answers
		FROM levels
		ORDER BY level_number ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var levels []dbModels.Level

	for rows.Next() {
		var level dbModels.Level

		err := rows.Scan(
			&level.ID,
			&level.LevelNumber,
			&level.Name,
			&level.WordsRequired,
			&level.MaxWrongAnswers,
		)
		if err != nil {
			return nil, err
		}

		levels = append(levels, level)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return levels, nil
}

func (r *LevelRepository) FindByID(id int) (*dbModels.Level, error) {
	query := `
		SELECT
			id,
			level_number,
			name,
			words_required,
			max_wrong_answers
			FROM levels
		WHERE id = ?
		LIMIT 1
	`

	return r.scanLevel(query, id)
}

func (r *LevelRepository) FindByLevelNumber(levelNumber int) (*dbModels.Level, error) {
	query := `
		SELECT
			id,
			level_number,
			name,
			words_required,
			max_wrong_answers
		FROM levels
		WHERE level_number = ?
		LIMIT 1
	`

	return r.scanLevel(query, levelNumber)
}

func (r *LevelRepository) GetNextLevel(currentLevelID int) (*dbModels.Level, error) {
	query := `
		SELECT
			next_level.id,
			next_level.level_number,
			next_level.name,
			next_level.words_required,
			next_level.max_wrong_answers
		FROM levels current_level
		JOIN levels next_level
			ON next_level.level_number = current_level.level_number + 1
		WHERE current_level.id = ?
		LIMIT 1
	`

	return r.scanLevel(query, currentLevelID)
}

func (r *LevelRepository) Exists(levelID int) (bool, error) {
	query := `
		SELECT COUNT(1)
		FROM levels
		WHERE id = ?
	`

	var count int

	err := r.db.QueryRow(query, levelID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *LevelRepository) scanLevel(query string, args ...any) (*dbModels.Level, error) {
	var level dbModels.Level

	err := r.db.QueryRow(query, args...).Scan(
		&level.ID,
		&level.LevelNumber,
		&level.Name,
		&level.WordsRequired,
		&level.MaxWrongAnswers,
	)

	if err != nil {
		return nil, err
	}

	return &level, nil
}
