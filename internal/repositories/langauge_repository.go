package repositories

import (
	"database/sql"

	"github.com/rajkovrga/1000words-game/internal/db"
	dbmodel "github.com/rajkovrga/1000words-game/internal/models/db"
)

type LanguageRepository struct {
	db *sql.DB
}

func NewLanguageRepository(db *sql.DB) LanguageRepository {
	return LanguageRepository{db: db}
}

func (languageRepository LanguageRepository) GetLanguagePerLanguageCode(code string) (*dbmodel.Language, error) {
	var language *dbmodel.Language

	err := db.DB.QueryRow(`
		SELECT id, code, language_value
		FROM languages
		WHERE code = ?
	`, code).Scan(
		&language.Id,
		&language.Code,
		&language.LanguageName,
	)

	if err != nil {
		return nil, err
	}

	return language, nil
}
