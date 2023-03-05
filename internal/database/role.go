package database

import (
	"context"
	"log"

	"github.com/segmentio/ksuid"
	"github.com/tapiaw38/auth-api/internal/models"
)

// EnsureRole ensures that the base roles are present
func (repository *PostgresRepository) EnsureRole() error {
	baseRoles := []string{"superadmin", "admin", "user", "guest"}
	for _, name := range baseRoles {
		r, _ := repository.GetRoleByName(context.Background(), name)

		if r != nil {
			continue
		}

		id, err := ksuid.NewRandom()

		if err != nil {
			return err
		}

		var role = models.Role{
			Id:   id.String(),
			Name: name,
		}

		_, err = repository.InsertRole(context.Background(), &role)

		if err != nil {
			return err
		}
	}

	return nil
}

// InsertRole inserts a role
func (repository *PostgresRepository) InsertRole(ctx context.Context, role *models.Role) (*models.Role, error) {

	q := `
		INSERT INTO roles (id, name)
		VALUES ($1, $2)
		RETURNING id, name
	`

	rows := repository.db.QueryRowContext(ctx, q, role.Id, role.Name)

	r, err := ScanRowRole(rows)
	if err != nil {
		return &models.Role{}, err
	}

	return &r, nil
}

// GetRoleByName returns a role by name
func (repository *PostgresRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {

	q := `
		SELECT id, name
		FROM roles
		WHERE name = $1
	`

	rows, err := repository.db.QueryContext(ctx, q, name)

	defer func() {
		err = rows.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	var role models.Role

	for rows.Next() {
		role, err = ScanRowRole(rows)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &role, nil
}

// GetRoleById returns a role by id
func (repository *PostgresRepository) GetRoleById(ctx context.Context, id string) (*models.Role, error) {

	q := `
		SELECT id, name
		FROM roles
		WHERE id = $1
	`

	rows, err := repository.db.QueryContext(ctx, q, id)

	defer func() {
		err = rows.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	var role models.Role

	for rows.Next() {
		role, err = ScanRowRole(rows)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &role, nil
}

// UpdateRole updates a role
func (repository *PostgresRepository) UpdateRole(ctx context.Context, role *models.Role) (*models.Role, error) {

	q := `
		UPDATE roles
		SET name = $1
		WHERE id = $2
		RETURNING id, name
	`

	rows := repository.db.QueryRowContext(ctx, q, role.Name, role.Id)

	r, err := ScanRowRole(rows)
	if err != nil {
		return &models.Role{}, err
	}

	return &r, nil
}

// DeleteRole deletes a role
func (repository *PostgresRepository) DeleteRole(ctx context.Context, id string) (*models.Role, error) {
	q := `
		DELETE FROM roles
		WHERE id = $1
		RETURNING id, name
	`

	row := repository.db.QueryRowContext(ctx, q, id)

	r, err := ScanRowRole(row)
	if err != nil {
		return &models.Role{}, err
	}

	return &r, nil
}

// ListRole returns a list of roles
func (repository *PostgresRepository) ListRole(ctx context.Context) ([]*models.Role, error) {
	q := `
		SELECT id, name
		FROM roles
	`

	rows, err := repository.db.QueryContext(ctx, q)

	defer func() {
		err = rows.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	var roles []*models.Role

	for rows.Next() {
		role, err := ScanRowRole(rows)
		if err != nil {
			return nil, err
		}

		roles = append(roles, &role)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}
