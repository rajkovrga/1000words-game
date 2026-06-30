package repositories

import "database/sql"

type AuthorizationRepository struct {
	db *sql.DB
}

func NewAuthorizationRepository(db *sql.DB) *AuthorizationRepository {
	return &AuthorizationRepository{
		db: db,
	}
}

func (r *AuthorizationRepository) AssignRoleToUserByCode(
	userID int,
	roleCode string,
) error {
	query := `
		INSERT OR IGNORE INTO user_roles (
			user_id,
			role_id
		)
		SELECT
			?,
			roles.id
		FROM roles
		WHERE roles.code = ?
		  AND roles.deleted_at IS NULL
	`

	_, err := r.db.Exec(query, userID, roleCode)

	return err
}

func (r *AuthorizationRepository) UserHasPermission(
	userID int,
	permissionCode string,
) (bool, error) {
	query := `
		SELECT COUNT(1)
		FROM user_roles
		JOIN roles
			ON roles.id = user_roles.role_id
		JOIN role_permissions
			ON role_permissions.role_id = roles.id
		JOIN permissions
			ON permissions.id = role_permissions.permission_id
		WHERE user_roles.user_id = ?
		  AND permissions.code = ?
		  AND roles.deleted_at IS NULL
		  AND permissions.deleted_at IS NULL
	`

	var count int

	err := r.db.QueryRow(
		query,
		userID,
		permissionCode,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AuthorizationRepository) GetUserRoleCodes(userID int) ([]string, error) {
	query := `
		SELECT DISTINCT
			roles.code
		FROM user_roles
		JOIN roles
			ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = ?
		  AND roles.deleted_at IS NULL
		ORDER BY roles.code ASC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string

	for rows.Next() {
		var roleCode string

		if err := rows.Scan(&roleCode); err != nil {
			return nil, err
		}

		roles = append(roles, roleCode)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *AuthorizationRepository) GetUserPermissionCodes(userID int) ([]string, error) {
	query := `
		SELECT DISTINCT
			permissions.code
		FROM user_roles
		JOIN roles
			ON roles.id = user_roles.role_id
		JOIN role_permissions
			ON role_permissions.role_id = roles.id
		JOIN permissions
			ON permissions.id = role_permissions.permission_id
		WHERE user_roles.user_id = ?
		  AND roles.deleted_at IS NULL
		  AND permissions.deleted_at IS NULL
		ORDER BY permissions.code ASC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string

	for rows.Next() {
		var permissionCode string

		if err := rows.Scan(&permissionCode); err != nil {
			return nil, err
		}

		permissions = append(permissions, permissionCode)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
