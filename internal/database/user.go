package database

import (
	"context"
	"log"
	"time"

	"github.com/tapiaw38/auth-api/internal/models"
)

// InsertUser inserts a new user into the database
func (repository *PostgresRepository) InsertUser(ctx context.Context, user *models.User) (*models.UserResponse, error) {
	q := `
		INSERT INTO users (
			id, first_name, last_name, username, email, 
			password, phone_number, picture, address, 
			is_active, verified_email, 
			created_at, updated_at
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, first_name, last_name, username, 
			email, phone_number, picture, address, 
			is_active, verified_email, verified_email_token,
			verified_email_token_expiry, password_reset_token,
			password_reset_token_expiry,
			created_at, updated_at;
		`
	row := repository.db.QueryRowContext(
		ctx, q,
		user.Id, user.FirstName, user.LastName,
		user.Username, user.Email, user.Password,
		user.PhoneNumber, user.Picture, user.Address,
		user.IsActive, user.VerifiedEmail,
		time.Now(), time.Now(),
	)

	u, err := ScanRowUserResponse(row)
	if err != nil {
		return &models.UserResponse{}, err
	}

	return u, nil
}

// GetUserById returns a user by id
func (repository *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.UserResponse, error) {

	q := `
		SELECT id, first_name, last_name, username, 
			email, phone_number, picture, address, 
			is_active, verified_email, verified_email_token,
			verified_email_token_expiry, password_reset_token,
			password_reset_token_expiry,
			created_at, updated_at
		FROM users 
		WHERE id = $1;
	`

	rows, err := repository.db.QueryContext(ctx, q, id)

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user *models.UserResponse

	for rows.Next() {
		user, err = ScanRowUserResponse(rows)
		if err != nil {
			return nil, err
		}

		q := `
			SELECT roles.id, roles.name
			FROM roles
			INNER JOIN user_roles
			ON roles.id = user_roles.role_id
			WHERE user_roles.user_id = $1;
		`

		rows, err = repository.db.QueryContext(ctx, q, user.Id)

		defer func() {
			err = rows.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		for rows.Next() {
			role, err := ScanRowRole(rows)
			if err != nil {
				return nil, err
			}

			user.Roles = append(user.Roles, *role)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByVerifiedEmailToken returns a user by verified email token
func (repository *PostgresRepository) GetUserByVerifiedEmailToken(ctx context.Context, token string) (*models.UserResponse, error) {

	q := `
		SELECT id, first_name, last_name, username, 
			email, phone_number, picture, address, 
			is_active, verified_email, verified_email_token,
			verified_email_token_expiry, password_reset_token,
			password_reset_token_expiry,
			created_at, updated_at
		FROM users 
		WHERE verified_email_token = $1;
	`

	rows, err := repository.db.QueryContext(ctx, q, token)

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user *models.UserResponse

	for rows.Next() {

		user, err = ScanRowUserResponse(rows)
		if err != nil {
			return nil, err
		}

		q := `
			SELECT roles.id, roles.name
			FROM roles
			INNER JOIN user_roles
			ON roles.id = user_roles.role_id
			WHERE user_roles.user_id = $1;
		`

		rows, err = repository.db.QueryContext(ctx, q, user.Id)

		defer func() {
			err = rows.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		for rows.Next() {
			role, err := ScanRowRole(rows)
			if err != nil {
				return nil, err
			}

			user.Roles = append(user.Roles, *role)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByVerifiedEmailToken returns a user by verified email token
func (repository *PostgresRepository) GetUserByPasswordResetToken(ctx context.Context, token string) (*models.UserResponse, error) {

	q := `
		SELECT id, first_name, last_name, username, 
			email, phone_number, picture, address, 
			is_active, verified_email, verified_email_token,
			verified_email_token_expiry, password_reset_token,
			password_reset_token_expiry,
			created_at, updated_at
		FROM users 
		WHERE password_reset_token = $1;
	`

	rows, err := repository.db.QueryContext(ctx, q, token)

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user *models.UserResponse

	for rows.Next() {

		user, err = ScanRowUserResponse(rows)
		if err != nil {
			return nil, err
		}

		q := `
			SELECT roles.id, roles.name
			FROM roles
			INNER JOIN user_roles
			ON roles.id = user_roles.role_id
			WHERE user_roles.user_id = $1;
		`

		rows, err = repository.db.QueryContext(ctx, q, user.Id)

		defer func() {
			err = rows.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		for rows.Next() {
			role, err := ScanRowRole(rows)
			if err != nil {
				return nil, err
			}

			user.Roles = append(user.Roles, *role)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return user, nil
}

func (repository *PostgresRepository) GetUserByEmailSocial(ctx context.Context, email string) (*models.UserResponse, error) {

	q := `
		SELECT id, first_name, last_name, username, 
			email, phone_number, picture, address, 
			is_active, verified_email, verified_email_token,
			verified_email_token_expiry, password_reset_token,
			password_reset_token_expiry,
			created_at, updated_at
		FROM users 
		WHERE email = $1;
	`

	rows, err := repository.db.QueryContext(ctx, q, email)

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user *models.UserResponse

	for rows.Next() {
		user, err = ScanRowUserResponse(rows)
		if err != nil {
			return nil, err
		}

		q := `
			SELECT roles.id, roles.name
			FROM roles
			INNER JOIN user_roles
			ON roles.id = user_roles.role_id
			WHERE user_roles.user_id = $1;
		`

		rows, err = repository.db.QueryContext(ctx, q, user.Id)

		defer func() {
			err = rows.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		for rows.Next() {
			role, err := ScanRowRole(rows)
			if err != nil {
				return nil, err
			}

			user.Roles = append(user.Roles, *role)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail returns a user by email
func (repository *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {

	q := `
		SELECT id, first_name, last_name, username, email, 
			password, phone_number, picture, address, 
			is_active, verified_email, verified_email_token, 
			verified_email_token_expiry, password_reset_token, 
			password_reset_token_expiry, created_at, updated_at
		FROM users 
		WHERE email = $1;
	`

	rows, err := repository.db.QueryContext(ctx, q, email)

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

		q := `
			SELECT roles.id, roles.name
			FROM roles
			INNER JOIN user_roles
			ON roles.id = user_roles.role_id
			WHERE user_roles.user_id = $1;
		`

		rows, err = repository.db.QueryContext(ctx, q, user.Id)

		defer func() {
			err = rows.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		for rows.Next() {
			role, err := ScanRowRole(rows)
			if err != nil {
				return nil, err
			}

			user.Roles = append(user.Roles, *role)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates a user in the database
func (ur *PostgresRepository) UpdateUser(ctx context.Context, id string, user *models.UserResponse) (*models.UserResponse, error) {
	q := `
	UPDATE users
		SET 
			first_name = $1, last_name = $2, email = $3, 
			picture = $4, phone_number = $5, address = $6, 
			is_active = $7, verified_email = $8, 
			verified_email_token = $9, verified_email_token_expiry = $10,
			password_reset_token = $11, password_reset_token_expiry = $12,
			updated_at = $13
		WHERE id = $14
		RETURNING id, first_name, last_name, username, 
			email, phone_number, picture, address, 
			is_active, verified_email, verified_email_token,
			verified_email_token_expiry, password_reset_token,
			password_reset_token_expiry,
			created_at, updated_at;
	`

	row := ur.db.QueryRowContext(
		ctx, q, user.FirstName, user.LastName, user.Email,
		user.Picture, user.PhoneNumber, user.Address,
		user.IsActive, user.VerifiedEmail, user.VerifiedEmailToken,
		user.VerifiedEmailTokenExpiry, user.PasswordResetToken,
		user.PasswordResetTokenExpiry,
		time.Now(), id,
	)

	u, err := ScanRowUserResponse(row)
	if err != nil {
		return &models.UserResponse{}, err
	}

	q = `
			SELECT roles.id, roles.name
			FROM roles
			INNER JOIN user_roles
			ON roles.id = user_roles.role_id
			WHERE user_roles.user_id = $1;
		`

	rows, err := ur.db.QueryContext(ctx, q, u.Id)

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		role, err := ScanRowRole(rows)
		if err != nil {
			return nil, err
		}

		u.Roles = append(u.Roles, *role)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return u, nil
}

// PartialUpdateUser updates a user in the database
func (ur *PostgresRepository) PartialUpdateUser(ctx context.Context, id string, user *models.UserResponse) (*models.UserResponse, error) {

	q := `
	UPDATE users
		SET 
			first_name = CASE WHEN $1 = '' THEN first_name ELSE $1 END, 
			last_name = CASE WHEN $2 = '' THEN last_name ELSE $2 END, 
			email = CASE WHEN $3 = '' THEN email ELSE $3 END,
			phone_number = CASE WHEN $4 = '' THEN phone_number ELSE $4 END,
			picture = CASE WHEN $5 = '' THEN picture ELSE $5 END,
			address = CASE WHEN $6 = '' THEN address ELSE $6 END,
			is_active = 
				CASE 
					WHEN $7 = TRUE AND is_active = FALSE THEN TRUE 
					WHEN $7 = FALSE AND is_active = TRUE THEN FALSE
					WHEN $7 = NULL THEN is_active
					ELSE is_active
				END,
			verified_email =
				CASE
					WHEN $8 = TRUE AND verified_email = FALSE THEN TRUE
					WHEN $8 = FALSE AND verified_email = TRUE THEN FALSE
					WHEN $8 = NULL THEN verified_email
					ELSE verified_email
				END,
			verified_email_token = CASE WHEN $9 = '' THEN verified_email_token ELSE $9 END,
			verified_email_token_expiry = $10,
			password_reset_token = CASE WHEN $11 = '' THEN password_reset_token ELSE $11 END,
			password_reset_token_expiry = $12,
			updated_at = $13
		WHERE id = $14
		RETURNING id, first_name, last_name, username, email, 
			phone_number, picture, address, 
			is_active, verified_email, verified_email_token,
			verified_email_token_expiry, password_reset_token,
			password_reset_token_expiry,
			created_at, updated_at
	`
	row := ur.db.QueryRowContext(
		ctx, q, user.FirstName, user.LastName, user.Email,
		user.PhoneNumber, user.Picture, user.Address,
		user.IsActive, user.VerifiedEmail, user.VerifiedEmailToken,
		user.VerifiedEmailTokenExpiry, user.PasswordResetToken,
		user.PasswordResetTokenExpiry,
		time.Now(), id,
	)

	u, err := ScanRowUserResponse(row)
	if err != nil {
		return &models.UserResponse{}, err
	}

	q = `
			SELECT roles.id, roles.name
			FROM roles
			INNER JOIN user_roles
			ON roles.id = user_roles.role_id
			WHERE user_roles.user_id = $1;
		`

	rows, err := ur.db.QueryContext(ctx, q, u.Id)

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		role, err := ScanRowRole(rows)
		if err != nil {
			return nil, err
		}

		u.Roles = append(u.Roles, *role)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return u, nil
}

// GetUsers returns all users
func (repository *PostgresRepository) ListUser(ctx context.Context) ([]*models.UserResponse, error) {

	q := `
	SELECT id, first_name, last_name, username, 
		email, phone_number, picture, address, 
		is_active, verified_email, verified_email_token,
		verified_email_token_expiry, password_reset_token,
		password_reset_token_expiry,
		created_at, updated_at
	FROM users
	`

	rows, err := repository.db.QueryContext(ctx, q)

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var users []*models.UserResponse

	for rows.Next() {

		user, err := ScanRowUserResponse(rows)
		if err != nil {
			return nil, err
		}

		q = `
			SELECT roles.id, roles.name
			FROM roles
			INNER JOIN user_roles
			ON roles.id = user_roles.role_id
			WHERE user_roles.user_id = $1;
		`

		rows, err := repository.db.QueryContext(ctx, q, user.Id)

		defer func() {
			err = rows.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		for rows.Next() {
			role, err := ScanRowRole(rows)
			if err != nil {
				return nil, err
			}

			user.Roles = append(user.Roles, *role)
		}

		if err = rows.Err(); err != nil {
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
