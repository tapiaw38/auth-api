package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"github.com/tapiaw38/auth-api/internal/models"
	"github.com/tapiaw38/auth-api/internal/repository"
	"github.com/tapiaw38/auth-api/internal/server"
	"github.com/tapiaw38/auth-api/internal/utils"
	"golang.org/x/crypto/bcrypt"
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

// SendEmailVerification sends an email verification email to a user
func SendEmailVerification(s server.Server, u *models.UserResponse, token string) error {
	// channel := make(chan error)
	subjet := "Bienvenido a Mi Tur"

	variables := map[string]string{
		"name": u.FirstName + " " + u.LastName,
		"link": s.Config().Host + "/auth/verify-email?token=" + token,
	}

	err := s.Rabbit().Connection().PublishEmailVerification(u.Email, s.Config().EmailHostUser, subjet, variables)
	if err != nil {
		return err
	}

	return nil
}

// HandleGoogleLogin handles the google login request
func HandleGoogleLogin(c *gin.Context, s server.Server, request *SignUpLoginRequest) (*models.UserResponse, error) {
	token, err := s.Google().ExchangeCode(c.Request.Context(), request.Code)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.Google().GetUserInfo(c.Request.Context(), token)
	if err != nil {
		return nil, err
	}

	user, err := repository.GetUserByEmailSocial(c.Request.Context(), userInfo.Email)
	if err != nil || user == nil {
		// If user not registered, register user
		id, err := ksuid.NewRandom()
		if err != nil {
			return nil, err
		}

		userInsert := models.User{
			Id:            id.String(),
			FirstName:     userInfo.FirstName,
			LastName:      userInfo.LastName,
			Username:      utils.RandomString(30),
			Email:         userInfo.Email,
			Password:      "",
			Picture:       userInfo.Picture,
			Address:       "",
			PhoneNumber:   "",
			IsActive:      true,
			VerifiedEmail: userInfo.VerifiedEmail,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		user, err = repository.InsertUser(c.Request.Context(), &userInsert)
		if err != nil {
			return nil, err
		}

		user, err := AddRoleToUser(c.Request.Context(), user.Id, "user")
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	// If user already registered, update user info
	if user.Picture == "" || !user.VerifiedEmail {
		userUpdate := models.UserResponse{
			Picture:       userInfo.Picture,
			IsActive:      user.IsActive,
			VerifiedEmail: userInfo.VerifiedEmail,
		}

		user, err = repository.PartialUpdateUser(c.Request.Context(), user.Id, &userUpdate)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return user, nil
}

// HandleEmailAndPasswordLogin handles the email and password login request
func HandleEmailAndPasswordLogin(c *gin.Context, s server.Server, request *SignUpLoginRequest) (*models.User, error) {

	user, err := repository.GetUserByEmail(c.Request.Context(), request.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// DecodeToten decodes a user token
func DecodeToken(tokenString, secret string) (*models.AppClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
