package repository

import (
	"context"

	"github.com/tapiaw38/auth-api/internal/models"
)

func InsertUser(ctx context.Context, user *models.User) (*models.User, error) {
	return implementation.InsertUser(ctx, user)
}

func GetUserById(ctx context.Context, id string) (*models.User, error) {
	return implementation.GetUserById(ctx, id)
}

func GetUserByVerifiedEmailToken(ctx context.Context, token string) (*models.User, error) {
	return implementation.GetUserByVerifiedEmailToken(ctx, token)
}

func GetUserByPasswordResetToken(ctx context.Context, token string) (*models.User, error) {
	return implementation.GetUserByPasswordResetToken(ctx, token)
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementation.GetUserByEmail(ctx, email)
}

func UpdateUser(ctx context.Context, id string, user *models.User) (*models.User, error) {
	return implementation.UpdateUser(ctx, id, user)
}

func UpdateUserProfile(ctx context.Context, id string, userProfile *models.UserProfile) (*models.User, error) {
	return implementation.UpdateUserProfile(ctx, id, userProfile)
}

func PartialUpdateUser(ctx context.Context, id string, user *models.User) (*models.User, error) {
	return implementation.PartialUpdateUser(ctx, id, user)
}

func ListUser(ctx context.Context, page int, limit int) ([]*models.User, error) {
	return implementation.ListUser(ctx, page, limit)
}

func Close() error {
	return implementation.Close()
}
