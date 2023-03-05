package repository

import (
	"context"

	"github.com/tapiaw38/auth-api/internal/models"
)

func EnsureRole() error {
	return implementation.EnsureRole()
}

func InsertRole(ctx context.Context, role *models.Role) (*models.Role, error) {
	return implementation.InsertRole(ctx, role)
}

func GetRoleById(ctx context.Context, id string) (*models.Role, error) {
	return implementation.GetRoleById(ctx, id)
}

func GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	return implementation.GetRoleByName(ctx, name)
}

func UpdateRole(ctx context.Context, role *models.Role) (*models.Role, error) {
	return implementation.UpdateRole(ctx, role)
}

func DeleteRole(ctx context.Context, id string) (*models.Role, error) {
	return implementation.DeleteRole(ctx, id)
}

func ListRole(ctx context.Context) ([]*models.Role, error) {
	return implementation.ListRole(ctx)
}
