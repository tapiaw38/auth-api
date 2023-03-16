package repository

import (
	"context"

	"github.com/tapiaw38/auth-api/internal/models"
)

// Repository is the repository interface
type Repository interface {
	// User
	InsertUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByVerifiedEmailToken(ctx context.Context, token string) (*models.User, error)
	GetUserByPasswordResetToken(ctx context.Context, token string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, user *models.User) (*models.User, error)
	PartialUpdateUser(ctx context.Context, id string, user *models.User) (*models.User, error)
	ListUser(ctx context.Context, page int, limit int) ([]*models.User, error)
	// Role
	EnsureRole() error
	InsertRole(ctx context.Context, role *models.Role) (*models.Role, error)
	GetRoleById(ctx context.Context, id string) (*models.Role, error)
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	UpdateRole(ctx context.Context, role *models.Role) (*models.Role, error)
	DeleteRole(ctx context.Context, id string) (*models.Role, error)
	ListRole(ctx context.Context) ([]*models.Role, error)
	// User Role
	InsertUserRole(ctx context.Context, userRole *models.UserRole) error
	DeleteUserRole(ctx context.Context, userRole *models.UserRole) error

	Close() error
}

var implementation Repository

func SetRepository(repository Repository) {
	implementation = repository
}
