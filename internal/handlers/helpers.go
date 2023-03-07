package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

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

// Generate a token for a user
func GenerateToken(user *models.UserResponse) (string, error) {
	// Create a random byte slice.
	tokenBytes := make([]byte, 16)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	// Encode the byte slice to a hexadecimal string.
	token := hex.EncodeToString(tokenBytes)

	// Save the token to the database.
	user.Token = token
	user.TokenExpiry = time.Now().Add(time.Hour * 48)

	userResponse := models.UserResponse{
		Token:         user.Token,
		TokenExpiry:   user.TokenExpiry,
		IsActive:      user.IsActive,
		VerifiedEmail: user.VerifiedEmail,
	}

	_, err = repository.PartialUpdateUser(context.Background(), user.Id, &userResponse)
	if err != nil {
		return "", err
	}

	return token, nil
}
