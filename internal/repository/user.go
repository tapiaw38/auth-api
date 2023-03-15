package repository

import (
	"context"

	"github.com/tapiaw38/auth-api/internal/models"
)

func InsertUser(ctx context.Context, user *models.User) (*models.UserResponse, error) {
	return implementation.InsertUser(ctx, user)
}

func GetUserById(ctx context.Context, id string) (*models.UserResponse, error) {
	return implementation.GetUserById(ctx, id)
}

func GetUserByVerifiedEmailToken(ctx context.Context, token string) (*models.UserResponse, error) {
	return implementation.GetUserByVerifiedEmailToken(ctx, token)
}

func GetUserByPasswordResetToken(ctx context.Context, token string) (*models.UserResponse, error) {
	return implementation.GetUserByPasswordResetToken(ctx, token)
}

func GetUserByEmailSocial(ctx context.Context, email string) (*models.UserResponse, error) {
	return implementation.GetUserByEmailSocial(ctx, email)
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementation.GetUserByEmail(ctx, email)
}

func UpdateUser(ctx context.Context, id string, user *models.UserResponse) (*models.UserResponse, error) {
	return implementation.UpdateUser(ctx, id, user)
}

func PartialUpdateUser(ctx context.Context, id string, user *models.UserResponse) (*models.UserResponse, error) {
	return implementation.PartialUpdateUser(ctx, id, user)
}

func ListUser(ctx context.Context) ([]*models.UserResponse, error) {
	return implementation.ListUser(ctx)
}

func Close() error {
	return implementation.Close()
}
