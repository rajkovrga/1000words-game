package repositories

import (
	"database/sql"
	"strings"

	dbModels "1000words-game/models/db"
	gameModels "1000words-game/models/game"
)

type WordRepository struct {
	db *sql.DB
}

func NewWordRepository(db *sql.DB) *WordRepository {
	return &WordRepository{
		db: db,
	}
}

func (r *WordRepository) GetQuestionsByLevel(
	levelID int,
	targetLanguageCode string,
	nativeLanguageCode string,
	limit int,
) ([]gameModels.Question, error) {
	query := `
		SELECT
			w.id AS word_id,
			target.text AS word,
			native.text AS answer
		FROM words w
		JOIN word_languages target
			ON target.word_id = w.id
		JOIN languages target_language
			ON target_language.id = target.language_id
		JOIN word_languages native
			ON native.word_id = w.id
		JOIN languages native_language
			ON native_language.id = native.language_id
		WHERE w.level_id = ?
		  AND target_language.code = ?
		  AND native_language.code = ?
		ORDER BY RANDOM()
		LIMIT ?
	`

	rows, err := r.db.Query(
		query,
		levelID,
		targetLanguageCode,
		nativeLanguageCode,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []gameModels.Question

	for rows.Next() {
		var question gameModels.Question

		err := rows.Scan(
			&question.WordID,
			&question.Word,
			&question.Answer,
		)
		if err != nil {
			return nil, err
		}

		questions = append(questions, question)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}

func (r *WordRepository) CreateWord(levelID int) (*dbModels.Word, error) {
	query := `
		INSERT INTO words (
			level_id
		)
		VALUES (?)
		RETURNING
			id,
			level_id,
			created_at
	`

	var word dbModels.Word

	err := r.db.QueryRow(query, levelID).Scan(
		&word.ID,
		&word.LevelID,
		&word.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &word, nil
}

func (r *WordRepository) AddWordLanguage(
	wordID int,
	languageID int,
	text string,
) (*dbModels.WordLanguage, error) {
	normalizedText := normalizeText(text)

	query := `
		INSERT INTO word_languages (
			word_id,
			language_id,
			text,
			normalized_text
		)
		VALUES (?, ?, ?, ?)
		RETURNING
			id,
			word_id,
			language_id,
			text,
			normalized_text,
			created_at
	`

	var wordLanguage dbModels.WordLanguage

	err := r.db.QueryRow(
		query,
		wordID,
		languageID,
		text,
		normalizedText,
	).Scan(
		&wordLanguage.ID,
		&wordLanguage.WordID,
		&wordLanguage.LanguageID,
		&wordLanguage.Text,
		&wordLanguage.NormalizedText,
		&wordLanguage.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &wordLanguage, nil
}

func (r *WordRepository) FindWordByID(id int) (*dbModels.Word, error) {
	query := `
		SELECT
			id,
			level_id,
			created_at
		FROM words
		WHERE id = ?
		LIMIT 1
	`

	var word dbModels.Word

	err := r.db.QueryRow(query, id).Scan(
		&word.ID,
		&word.LevelID,
		&word.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &word, nil
}

func (r *WordRepository) GetWordLanguages(wordID int) ([]dbModels.WordLanguage, error) {
	query := `
		SELECT
			id,
			word_id,
			language_id,
			text,
			normalized_text,
			created_at
		FROM word_languages
		WHERE word_id = ?
		ORDER BY language_id ASC
	`

	rows, err := r.db.Query(query, wordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wordLanguages []dbModels.WordLanguage

	for rows.Next() {
		var wordLanguage dbModels.WordLanguage

		err := rows.Scan(
			&wordLanguage.ID,
			&wordLanguage.WordID,
			&wordLanguage.LanguageID,
			&wordLanguage.Text,
			&wordLanguage.NormalizedText,
			&wordLanguage.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		wordLanguages = append(wordLanguages, wordLanguage)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return wordLanguages, nil
}

func normalizeText(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)

	return value
}
