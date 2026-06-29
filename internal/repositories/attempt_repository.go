package repositories

import (
	"database/sql"

	dbModels "1000words-game/models/db"
	gameModels "1000words-game/models/game"
)

type AttemptRepository struct {
	db *sql.DB
}

func NewAttemptRepository(db *sql.DB) *AttemptRepository {
	return &AttemptRepository{
		db: db,
	}
}

func (r *AttemptRepository) Create(
	userID int,
	targetLanguageID int,
	nativeLanguageID int,
	levelID int,
	totalQuestions int,
) (*dbModels.UserAttempt, error) {
	query := `
		INSERT INTO user_attempts (
			user_id,
			target_language_id,
			native_language_id,
			level_id,
			total_questions
		)
		VALUES (?, ?, ?, ?, ?)
		RETURNING
			id,
			user_id,
			target_language_id,
			native_language_id,
			level_id,
			correct_count,
			wrong_count,
			total_questions,
			passed,
			started_at,
			finished_at
	`

	return r.scanAttempt(
		query,
		userID,
		targetLanguageID,
		nativeLanguageID,
		levelID,
		totalQuestions,
	)
}

func (r *AttemptRepository) Finish(
	attemptID int,
	result gameModels.Result,
) (*dbModels.UserAttempt, error) {
	query := `
		UPDATE user_attempts
		SET
			correct_count = ?,
			wrong_count = ?,
			total_questions = ?,
			passed = ?,
			finished_at = CURRENT_TIMESTAMP
		WHERE id = ?
		RETURNING
			id,
			user_id,
			target_language_id,
			native_language_id,
			level_id,
			correct_count,
			wrong_count,
			total_questions,
			passed,
			started_at,
			finished_at
	`

	passed := 0
	if result.Passed {
		passed = 1
	}

	return r.scanAttempt(
		query,
		result.CorrectCount,
		result.WrongCount,
		result.TotalQuestions,
		passed,
		attemptID,
	)
}

func (r *AttemptRepository) FindByID(id int) (*dbModels.UserAttempt, error) {
	query := `
		SELECT
			id,
			user_id,
			target_language_id,
			native_language_id,
			level_id,
			correct_count,
			wrong_count,
			total_questions,
			passed,
			started_at,
			finished_at
		FROM user_attempts
		WHERE id = ?
		LIMIT 1
	`

	return r.scanAttempt(query, id)
}

func (r *AttemptRepository) GetUserAttempts(
	userID int,
	targetLanguageID int,
	nativeLanguageID int,
) ([]dbModels.UserAttempt, error) {
	query := `
		SELECT
			id,
			user_id,
			target_language_id,
			native_language_id,
			level_id,
			correct_count,
			wrong_count,
			total_questions,
			passed,
			started_at,
			finished_at
		FROM user_attempts
		WHERE user_id = ?
		  AND target_language_id = ?
		  AND native_language_id = ?
		ORDER BY started_at DESC
	`

	rows, err := r.db.Query(query, userID, targetLanguageID, nativeLanguageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []dbModels.UserAttempt

	for rows.Next() {
		var attempt dbModels.UserAttempt
		var finishedAt sql.NullTime
		var passed int

		err := rows.Scan(
			&attempt.ID,
			&attempt.UserID,
			&attempt.TargetLanguageID,
			&attempt.NativeLanguageID,
			&attempt.LevelID,
			&attempt.CorrectCount,
			&attempt.WrongCount,
			&attempt.TotalQuestions,
			&passed,
			&attempt.StartedAt,
			&finishedAt,
		)
		if err != nil {
			return nil, err
		}

		attempt.Passed = passed == 1

		if finishedAt.Valid {
			attempt.FinishedAt = &finishedAt.Time
		}

		attempts = append(attempts, attempt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return attempts, nil
}

func (r *AttemptRepository) GetLastUserAttempt(
	userID int,
	targetLanguageID int,
	nativeLanguageID int,
	levelID int,
) (*dbModels.UserAttempt, error) {
	query := `
		SELECT
			id,
			user_id,
			target_language_id,
			native_language_id,
			level_id,
			correct_count,
			wrong_count,
			total_questions,
			passed,
			started_at,
			finished_at
		FROM user_attempts
		WHERE user_id = ?
		  AND target_language_id = ?
		  AND native_language_id = ?
		  AND level_id = ?
		ORDER BY started_at DESC
		LIMIT 1
	`

	return r.scanAttempt(query, userID, targetLanguageID, nativeLanguageID, levelID)
}

func (r *AttemptRepository) scanAttempt(query string, args ...any) (*dbModels.UserAttempt, error) {
	var attempt dbModels.UserAttempt
	var finishedAt sql.NullTime
	var passed int

	err := r.db.QueryRow(query, args...).Scan(
		&attempt.ID,
		&attempt.UserID,
		&attempt.TargetLanguageID,
		&attempt.NativeLanguageID,
		&attempt.LevelID,
		&attempt.CorrectCount,
		&attempt.WrongCount,
		&attempt.TotalQuestions,
		&passed,
		&attempt.StartedAt,
		&finishedAt,
	)

	if err != nil {
		return nil, err
	}

	attempt.Passed = passed == 1

	if finishedAt.Valid {
		attempt.FinishedAt = &finishedAt.Time
	}

	return &attempt, nil
}
