package repository

import (
	"context"

	"github.com/tapiaw38/auth-api/models"
)

func InsertUserRole(ctx context.Context, userRole *models.UserRole) error {
	return implementation.InsertUserRole(ctx, userRole)
}

func DeleteUserRole(ctx context.Context, userRole *models.UserRole) error {
	return implementation.DeleteUserRole(ctx, userRole)
}
