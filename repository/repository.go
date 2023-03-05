package repository

import (
	"context"

	"github.com/tapiaw38/auth-api/models"
)

// Repository is the repository interface
type Repository interface {
	// User
	InsertUser(ctx context.Context, user *models.User) (*models.UserResponse, error)
	GetUserById(ctx context.Context, id string) (*models.UserResponse, error)
	GetUserByEmailSocial(ctx context.Context, email string) (*models.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, user *models.UserResponse) (*models.UserResponse, error)
	PartialUpdateUser(ctx context.Context, id string, user *models.UserResponse) (*models.UserResponse, error)
	ListUser(ctx context.Context) ([]*models.UserResponse, error)
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
