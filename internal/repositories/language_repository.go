package repositories

import (
	"database/sql"

	dbModels "1000words-game/models/db"
)

type LanguageRepository struct {
	db *sql.DB
}

func NewLanguageRepository(db *sql.DB) *LanguageRepository {
	return &LanguageRepository{
		db: db,
	}
}

func (r *LanguageRepository) GetAll() ([]dbModels.Language, error) {
	query := `
		SELECT
			id,
			name,
			code
					FROM languages
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var languages []dbModels.Language

	for rows.Next() {
		var language dbModels.Language

		err := rows.Scan(
			&language.ID,
			&language.Name,
			&language.Code,
		)
		if err != nil {
			return nil, err
		}

		languages = append(languages, language)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return languages, nil
}

func (r *LanguageRepository) FindByID(id int) (*dbModels.Language, error) {
	query := `
		SELECT
			id,
			name,
			code
		FROM languages
		WHERE id = ?
		LIMIT 1
	`

	return r.scanLanguage(query, id)
}

func (r *LanguageRepository) FindByCode(code string) (*dbModels.Language, error) {
	query := `
		SELECT
			id,
			name,
			code
		FROM languages
		WHERE code = ?
		LIMIT 1
	`

	return r.scanLanguage(query, code)
}

func (r *LanguageRepository) ExistsByID(id int) (bool, error) {
	query := `
		SELECT COUNT(1)
		FROM languages
		WHERE id = ?
	`

	var count int

	err := r.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *LanguageRepository) ExistsByCode(code string) (bool, error) {
	query := `
		SELECT COUNT(1)
		FROM languages
		WHERE code = ?
	`

	var count int

	err := r.db.QueryRow(query, code).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *LanguageRepository) scanLanguage(query string, args ...any) (*dbModels.Language, error) {
	var language dbModels.Language

	err := r.db.QueryRow(query, args...).Scan(
		&language.ID,
		&language.Name,
		&language.Code,
	)

	if err != nil {
		return nil, err
	}

	return &language, nil
}
