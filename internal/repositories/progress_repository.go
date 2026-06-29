package repositories

import (
	"database/sql"

	dbModels "1000words-game/models/db"
	gameModels "1000words-game/models/game"
)

type ProgressRepository struct {
	db *sql.DB
}

func NewProgressRepository(db *sql.DB) *ProgressRepository {
	return &ProgressRepository{
		db: db,
	}
}

func (r *ProgressRepository) Create(
	userID int,
	targetLanguageID int,
	nativeLanguageID int,
	currentLevelID int,
) (*dbModels.UserLanguageProgress, error) {
	query := `
		INSERT INTO user_language_progress (
			user_id,
			target_language_id,
			native_language_id,
			current_level_id
		)
		VALUES (?, ?, ?, ?)
		RETURNING
			id,
			user_id,
			target_language_id,
			native_language_id,
			current_level_id,
			total_attempts,
			total_passed,
			total_failed,
			total_correct_answers,
			total_wrong_answers,
			updated_at,
			created_at
	`

	return r.scanProgress(
		query,
		userID,
		targetLanguageID,
		nativeLanguageID,
		currentLevelID,
	)
}

func (r *ProgressRepository) GetOptionsByUserID(userID int) ([]gameModels.ProgressOption, error) {
	query := `
		SELECT
			p.id,
			p.user_id,
			p.target_language_id,
			p.native_language_id,
			p.current_level_id,
			l.level_number,
			target.code,
			target.name,
			native.code,
			native.name
		FROM user_language_progress p
		JOIN levels l
			ON l.id = p.current_level_id
		JOIN languages target
			ON target.id = p.target_language_id
		JOIN languages native
			ON native.id = p.native_language_id
		WHERE p.user_id = ?
		ORDER BY p.id ASC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []gameModels.ProgressOption

	for rows.Next() {
		var option gameModels.ProgressOption

		err := rows.Scan(
			&option.ProgressID,
			&option.UserID,
			&option.TargetLanguageID,
			&option.NativeLanguageID,
			&option.CurrentLevelID,
			&option.LevelNumber,
			&option.TargetCode,
			&option.TargetName,
			&option.NativeCode,
			&option.NativeName,
		)
		if err != nil {
			return nil, err
		}

		options = append(options, option)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return options, nil
}

func (r *ProgressRepository) FindOptionByID(progressID int) (*gameModels.ProgressOption, error) {
	query := `
		SELECT
			p.id,
			p.user_id,
			p.target_language_id,
			p.native_language_id,
			p.current_level_id,
			l.level_number,
			target.code,
			target.name,
			native.code,
			native.name
		FROM user_language_progress p
		JOIN levels l
			ON l.id = p.current_level_id
		JOIN languages target
			ON target.id = p.target_language_id
		JOIN languages native
			ON native.id = p.native_language_id
		WHERE p.id = ?
		LIMIT 1
	`

	var option gameModels.ProgressOption

	err := r.db.QueryRow(query, progressID).Scan(
		&option.ProgressID,
		&option.UserID,
		&option.TargetLanguageID,
		&option.NativeLanguageID,
		&option.CurrentLevelID,
		&option.LevelNumber,
		&option.TargetCode,
		&option.TargetName,
		&option.NativeCode,
		&option.NativeName,
	)

	if err != nil {
		return nil, err
	}

	return &option, nil
}

func (r *ProgressRepository) UpdateAfterAttempt(
	progressID int,
	result gameModels.Result,
	nextLevelID int,
) error {
	passed := 0
	failed := 1

	if result.Passed {
		passed = 1
		failed = 0
	}

	query := `
		UPDATE user_language_progress
		SET
			total_attempts = total_attempts + 1,
			total_passed = total_passed + ?,
			total_failed = total_failed + ?,
			total_correct_answers = total_correct_answers + ?,
			total_wrong_answers = total_wrong_answers + ?,
			current_level_id = CASE 
				WHEN ? > 0 THEN ?
				ELSE current_level_id
			END,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		passed,
		failed,
		result.CorrectCount,
		result.WrongCount,
		nextLevelID,
		nextLevelID,
		progressID,
	)

	return err
}

func (r *ProgressRepository) scanProgress(query string, args ...any) (*dbModels.UserLanguageProgress, error) {
	var progress dbModels.UserLanguageProgress

	err := r.db.QueryRow(query, args...).Scan(
		&progress.ID,
		&progress.UserID,
		&progress.TargetLanguageID,
		&progress.NativeLanguageID,
		&progress.CurrentLevelID,
		&progress.TotalAttempts,
		&progress.TotalPassed,
		&progress.TotalFailed,
		&progress.TotalCorrectAnswers,
		&progress.TotalWrongAnswers,
		&progress.UpdatedAt,
		&progress.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &progress, nil
}
