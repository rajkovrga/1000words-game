package repositories

import (
	"database/sql"

	"github.com/rajkovrga/1000words-game/internal/db"
	dbmodel "github.com/rajkovrga/1000words-game/internal/models/db"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (userRepository *UserRepository) GetFirstUser() (*dbmodel.User, error) {
	var user dbmodel.User

	err := db.DB.QueryRow(`
	SELECT id, email, native_langauge_id, role_id, created_at, updated_at
	FROM users
	LIMIT 1
	`).Scan(
		&user.Id,
		&user.Email,
		&user.NativeLanguageId,
		&user.RoleId,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
