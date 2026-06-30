package repositories

import (
	"database/sql"
	"time"

	dbModels "1000words-game/models/db"
)

type AuthTokenRepository struct {
	db *sql.DB
}

func NewAuthTokenRepository(db *sql.DB) *AuthTokenRepository {
	return &AuthTokenRepository{
		db: db,
	}
}

func (r *AuthTokenRepository) Create(
	userID int,
	tokenHash string,
	expiresAt time.Time,
) (*dbModels.AuthToken, error) {
	query := `
		INSERT INTO auth_tokens (
			user_id,
			token_hash,
			expires_at
		)
		VALUES (?, ?, ?)
		RETURNING
			id,
			user_id,
			token_hash,
			expires_at,
			created_at
	`

	var token dbModels.AuthToken

	err := r.db.QueryRow(
		query,
		userID,
		tokenHash,
		expiresAt,
	).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *AuthTokenRepository) FindValidByHash(
	tokenHash string,
	now time.Time,
) (*dbModels.AuthToken, error) {
	query := `
		SELECT
			id,
			user_id,
			token_hash,
			expires_at,
			created_at
		FROM auth_tokens
		WHERE token_hash = ?
		  AND expires_at > ?
		LIMIT 1
	`

	var token dbModels.AuthToken

	err := r.db.QueryRow(
		query,
		tokenHash,
		now,
	).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *AuthTokenRepository) DeleteByHash(tokenHash string) error {
	query := `
		DELETE FROM auth_tokens
		WHERE token_hash = ?
	`

	_, err := r.db.Exec(query, tokenHash)

	return err
}

func (r *AuthTokenRepository) DeleteExpired(now time.Time) error {
	query := `
		DELETE FROM auth_tokens
		WHERE expires_at <= ?
	`

	_, err := r.db.Exec(query, now)

	return err
}
