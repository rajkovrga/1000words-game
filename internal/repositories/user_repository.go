package repositories

import (
	"database/sql"

	dbModels "1000words-game/models/db"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(email string, password string) (*dbModels.User, error) {
	query := `
		INSERT INTO users (
			email,
			password
		)
		VALUES (?, ?)
		RETURNING 
			id,
			email,
			password,
			created_at,
			updated_at,
			deleted_at
	`

	var user dbModels.User
	var deletedAt sql.NullTime

	err := r.db.QueryRow(query, email, password).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		return nil, err
	}

	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	return &user, nil
}

func (r *UserRepository) FindByID(id int) (*dbModels.User, error) {
	query := `
		SELECT 
			id,
			email,
			password,
			created_at,
			updated_at,
			deleted_at
		FROM users
		WHERE id = ?
		  AND deleted_at IS NULL
		LIMIT 1
	`

	return r.scanUser(query, id)
}

func (r *UserRepository) FindByEmail(email string) (*dbModels.User, error) {
	query := `
		SELECT 
			id,
			email,
			password,
			created_at,
			updated_at,
			deleted_at
		FROM users
		WHERE email = ?
		  AND deleted_at IS NULL
		LIMIT 1
	`

	return r.scanUser(query, email)
}

func (r *UserRepository) UpdatePassword(userID int, password string) error {
	query := `
		UPDATE users
		SET password = ?,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		  AND deleted_at IS NULL
	`

	_, err := r.db.Exec(query, password, userID)
	return err
}

func (r *UserRepository) SoftDelete(userID int) error {
	query := `
		UPDATE users
		SET deleted_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		  AND deleted_at IS NULL
	`

	_, err := r.db.Exec(query, userID)
	return err
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	query := `
		SELECT COUNT(1)
		FROM users
		WHERE email = ?
		  AND deleted_at IS NULL
	`

	var count int

	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *UserRepository) scanUser(query string, args ...any) (*dbModels.User, error) {
	var user dbModels.User
	var deletedAt sql.NullTime

	err := r.db.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		return nil, err
	}

	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	return &user, nil
}
