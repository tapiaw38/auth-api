package database

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/tapiaw38/auth-api/internal/models"
)

// InsertUser inserts a new user into the database
func (repository *PostgresRepository) InsertUser(ctx context.Context, user *models.User) (*models.User, error) {
	q := `
		INSERT INTO users (
			id, first_name, last_name, username, email, 
			password, phone_number, picture, address,
			is_active, verified_email, 
			created_at, updated_at
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, first_name, last_name, username, 
			email, password, phone_number, picture, address, 
			is_active, verified_email, verified_email_token,
			verified_email_token_expiry, password_reset_token,
			password_reset_token_expiry, created_at, updated_at;
		`
	row := repository.db.QueryRowContext(
		ctx, q,
		user.Id, user.FirstName, user.LastName,
		user.Username, user.Email, user.Password,
		user.PhoneNumber, user.Picture, user.Address,
		user.IsActive, user.VerifiedEmail,
		time.Now(), time.Now(),
	)

	u, err := ScanRowUser(row)
	if err != nil {
		return &models.User{}, err
	}

	return u, nil
}

// updateUserRoles updates the roles of a user
func (ur *PostgresRepository) updateUserRoles(ctx context.Context, u *models.User) error {
	q := `
		SELECT roles.id, roles.name
		FROM roles
		INNER JOIN user_roles
		ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = $1;
	`

	rows, err := ur.db.QueryContext(ctx, q, u.Id)
	if err != nil {
		return err
	}
	defer rows.Close()

	u.Roles = nil // Limpiamos los roles actuales antes de actualizarlos

	for rows.Next() {
		role, err := ScanRowRole(rows)
		if err != nil {
			return err
		}
		u.Roles = append(u.Roles, *role)
	}

	return rows.Err()
}

// getUserByQuery returns a user by executing the given query
func (repository *PostgresRepository) getUserByQuery(ctx context.Context, query string, args ...interface{}) (*models.User, error) {
	rows, err := repository.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user *models.User

	for rows.Next() {
		user, err = ScanRowUser(rows)
		if err != nil {
			return nil, err
		}

		if err = repository.updateUserRoles(ctx, user); err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserById returns a user by id
func (repository *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, username, 
			email, password, phone_number, picture, address, 
			is_active, verified_email, verified_email_token,
			verified_email_token_expiry, password_reset_token,
			password_reset_token_expiry,
			created_at, updated_at
		FROM users 
		WHERE id = $1;
	`

	return repository.getUserByQuery(ctx, query, id)
}

// GetUserByVerifiedEmailToken returns a user by verified email token
func (repository *PostgresRepository) GetUserByVerifiedEmailToken(ctx context.Context, token string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, username, 
			email, password, phone_number, picture, address, 
			is_active, verified_email, verified_email_token,
			verified_email_token_expiry, password_reset_token,
			password_reset_token_expiry,
			created_at, updated_at
		FROM users 
		WHERE verified_email_token = $1;
	`

	return repository.getUserByQuery(ctx, query, token)
}

// GetUserByPasswordResetToken returns a user by password reset token
func (repository *PostgresRepository) GetUserByPasswordResetToken(ctx context.Context, token string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, username, email, 
			password, phone_number, picture, address, 
			is_active, verified_email, verified_email_token, 
			verified_email_token_expiry, password_reset_token, 
			password_reset_token_expiry, created_at, updated_at
		FROM users 
		WHERE password_reset_token = $1;
	`

	return repository.getUserByQuery(ctx, query, token)
}

// GetUserByEmail returns a user by email
func (repository *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, username, email, 
			password, phone_number, picture, address, 
			is_active, verified_email, verified_email_token, 
			verified_email_token_expiry, password_reset_token, 
			password_reset_token_expiry, created_at, updated_at
		FROM users 
		WHERE email = $1;
	`

	return repository.getUserByQuery(ctx, query, email)
}

// UpdateUser updates a user in the database
func (ur *PostgresRepository) UpdateUser(ctx context.Context, id string, user *models.User) (*models.User, error) {
	q := `
		UPDATE users
		SET 
			first_name = $1, last_name = $2, email = $3,
			password = $4, picture = $5, phone_number = $6, 
			address = $7, is_active = $8, verified_email = $9, 
			verified_email_token = $10, verified_email_token_expiry = $11,
			password_reset_token = $12, password_reset_token_expiry = $13,
			updated_at = $14
		WHERE id = $15
		RETURNING id, first_name, last_name, username, 
			email, password, phone_number, picture, address, 
			is_active, verified_email, verified_email_token,
			verified_email_token_expiry, password_reset_token,
			password_reset_token_expiry,
			created_at, updated_at;
	`

	row := ur.db.QueryRowContext(
		ctx, q, user.FirstName, user.LastName, user.Email,
		user.Password, user.Picture, user.PhoneNumber, user.Address,
		user.IsActive, user.VerifiedEmail, user.VerifiedEmailToken,
		user.VerifiedEmailTokenExpiry, user.PasswordResetToken,
		user.PasswordResetTokenExpiry,
		time.Now(), id,
	)

	u, err := ScanRowUser(row)
	if err != nil {
		return nil, err
	}

	if err := ur.updateUserRoles(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

// PartialUpdateUser partially updates a user in the database
func (ur *PostgresRepository) PartialUpdateUser(ctx context.Context, id string, updates map[string]interface{}) (*models.User, error) {
	// Comienza una transacción
	tx, err := ur.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Construye la consulta de actualización dinámicamente
	updateFields := []string{}
	values := []interface{}{}

	for field, value := range updates {
		if value != nil {
			updateFields = append(updateFields, field+" = $"+strconv.Itoa(len(values)+1))
			values = append(values, value)
		}
	}

	// Agregar el id al final de los valores
	values = append(values, id)

	if len(updateFields) == 0 {
		// No se realizaron actualizaciones
		return nil, nil
	}

	q := `
		UPDATE users
		SET ` + strings.Join(updateFields, ", ") + `
		WHERE id = $` + strconv.Itoa(len(values)) + `
		RETURNING id, first_name, last_name, username, email,
			password, phone_number, picture, address,
			is_active, verified_email, verified_email_token,
			verified_email_token_expiry, password_reset_token,
			password_reset_token_expiry,
			created_at, updated_at
	`

	row := tx.QueryRowContext(ctx, q, values...)

	u, err := ScanRowUser(row)
	if err != nil {
		return nil, err
	}

	if err := ur.updateUserRoles(ctx, u); err != nil {
		return nil, err
	}

	// Commit la transacción
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return u, nil
}

// GetUsers returns all users
func (repository *PostgresRepository) ListUser(ctx context.Context, limit int, page int) ([]*models.User, error) {

	q := `
	SELECT id, first_name, last_name, username, 
		email, password, phone_number, picture, address, 
		is_active, verified_email, verified_email_token,
		verified_email_token_expiry, password_reset_token,
		password_reset_token_expiry,
		created_at, updated_at
	FROM users
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2;
	`

	rows, err := repository.db.QueryContext(ctx, q, limit, (page-1)*limit)

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var users []*models.User

	for rows.Next() {

		user, err := ScanRowUser(rows)
		if err != nil {
			return nil, err
		}

		err = repository.updateUserRoles(ctx, user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Close closes the database connection
func (repository *PostgresRepository) Close() error {
	return repository.db.Close()
}
