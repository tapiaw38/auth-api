package handlers

import (
	"context"

	"github.com/tapiaw38/auth-api/internal/models"
	"github.com/tapiaw38/auth-api/internal/repository"
)

// Add a role to a user by name
func AddRoleToUser(ctx context.Context, userId, roleName string) (*models.UserResponse, error) {

	role, err := repository.GetRoleByName(ctx, roleName)
	if err != nil {
		return nil, err
	}

	userRole := models.UserRole{
		UserId: userId,
		RoleId: role.Id,
	}

	err = repository.InsertUserRole(ctx, &userRole)
	if err != nil {
		return nil, err
	}

	user, err := repository.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
