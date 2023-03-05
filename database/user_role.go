package database

import (
	"context"

	"github.com/tapiaw38/auth-api/models"
)

// InsertUserRole inserts a new user_role into the database
func (repository *PostgresRepository) InsertUserRole(ctx context.Context, userRole *models.UserRole) error {

	q := `
		INSERT INTO user_roles (
			user_id, role_id
		) VALUES ($1, $2)
	`

	_, err := repository.db.ExecContext(ctx, q, userRole.UserId, userRole.RoleId)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserRole deletes a user_role from the database
func (repository *PostgresRepository) DeleteUserRole(ctx context.Context, userRole *models.UserRole) error {

	q := `
		DELETE FROM user_roles
		WHERE user_id = $1 AND role_id = $2
		`

	_, err := repository.db.ExecContext(ctx, q, userRole.UserId, userRole.RoleId)
	if err != nil {
		return err
	}

	return nil
}
