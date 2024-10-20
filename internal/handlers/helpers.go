package handlers

import (
	"context"
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
func AddRoleToUser(ctx context.Context, userId, roleName string) (*models.User, error) {

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

// SaveVerifiedEmailToken saves the verified email token to the database
func SaveVerifiedEmailToken(ctx context.Context, user *models.User, token string) error {
	// Save the token to the database.
	user.VerifiedEmailToken = token
	user.VerifiedEmailTokenExpiry = time.Now().Add(time.Hour * 168)

	updates := map[string]interface{}{
		"verified_email_token":        user.VerifiedEmailToken,
		"verified_email_token_expiry": user.VerifiedEmailTokenExpiry,
		"updated_at":                  time.Now(),
	}

	_, err := repository.PartialUpdateUser(ctx, user.Id, updates)
	if err != nil {
		return err
	}

	return nil
}

// SavePasswordResetToken saves the password reset token to the database
func SavePasswordResetToken(ctx context.Context, user *models.User, token string) error {
	// Save the token to the database.
	user.PasswordResetToken = token
	user.PasswordResetTokenExpiry = time.Now().Add(time.Hour * 2)

	updates := map[string]interface{}{
		"password_reset_token":        user.PasswordResetToken,
		"password_reset_token_expiry": user.PasswordResetTokenExpiry,
		"updated_at":                  time.Now(),
	}

	_, err := repository.PartialUpdateUser(ctx, user.Id, updates)
	if err != nil {
		return err
	}

	return nil
}

// DestroyVerifiedEmailToken destroys the verified email token
func DestroyPasswordResetToken(ctx context.Context, user *models.User) error {
	// Save the token to the database.
	user.PasswordResetToken = ""
	user.PasswordResetTokenExpiry = time.Time{}

	updates := map[string]interface{}{
		"password_reset_token":        user.PasswordResetToken,
		"password_reset_token_expiry": user.PasswordResetTokenExpiry,
		"updated_at":                  time.Now(),
	}

	_, err := repository.PartialUpdateUser(ctx, user.Id, updates)
	if err != nil {
		return err
	}

	return nil
}

// SendVerificationEmail sends an email verification email to a user
func SendVerificationEmail(s server.Server, u *models.User, token string) error {

	templateName := "email_verification"
	subjet := "Bienvenido a Mi Tur"

	variables := map[string]string{
		"name": u.FirstName + " " + u.LastName,
		"link": s.Config().Host + "/auth/verify-email?token=" + token,
	}

	err := s.Rabbit().Connection().PublishEmailMessage(u.Email, s.Config().EmailHostUser, subjet, templateName, variables)
	if err != nil {
		return err
	}

	return nil
}

// SendResetPasswordEmail sends an email verification email to a user
func SendResetPasswordEmail(s server.Server, u *models.User, token string) error {

	templateName := "reset_password"
	subjet := "Restablecer contraseña"

	variables := map[string]string{
		"name": u.FirstName + " " + u.LastName,
		"link": s.Config().FrontendURL + "/auth/reset-password?token=" + token,
	}

	err := s.Rabbit().Connection().PublishEmailMessage(u.Email, s.Config().EmailHostUser, subjet, templateName, variables)
	if err != nil {
		return err
	}

	return nil
}

// HandleGoogleLogin handles the google login request
func HandleGoogleLogin(c *gin.Context, s server.Server, request *SignUpLoginRequest) (*models.User, error) {
	token, err := s.Google().ExchangeCode(c.Request.Context(), request.Code)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.Google().GetUserInfo(c.Request.Context(), token)
	if err != nil {
		return nil, err
	}

	user, err := repository.GetUserByEmail(c.Request.Context(), userInfo.Email)
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
		userUpdate := models.User{
			Picture:       userInfo.Picture,
			VerifiedEmail: userInfo.VerifiedEmail,
		}

		updates := map[string]interface{}{
			"picture":        userUpdate.Picture,
			"verified_email": userUpdate.VerifiedEmail,
			"updated_at":     time.Now(),
		}

		user, err = repository.PartialUpdateUser(c.Request.Context(), user.Id, updates)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return user, nil
}

// HandleEmailAndPasswordLogin handles the email and password login request
func HandleEmailAndPasswordLogin(c *gin.Context, request *SignUpLoginRequest) (*models.User, error) {

	user, err := repository.GetUserByEmail(c.Request.Context(), request.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if err = ComparePassword(request.Password, user.Password); err != nil {
		return nil, err
	}

	return user, nil
}

// ComparePassword compares a password with a hash
func ComparePassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("invalid credentials")
	}

	return nil
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

// GetUserResponse returns a user without password
func GetUserResponse(user *models.User) *models.UserResponse {
	return &models.UserResponse{
		Id:            user.Id,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Username:      user.Username,
		Email:         user.Email,
		PhoneNumber:   user.PhoneNumber,
		Picture:       user.Picture,
		Address:       user.Address,
		IsActive:      user.IsActive,
		VerifiedEmail: user.VerifiedEmail,
		Roles:         user.Roles,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}

// GetUsersResponse returns a list of users without password
func GetUsersResponse(users []*models.User) []*models.UserResponse {
	var usersWithoutPassword []*models.UserResponse
	for _, user := range users {
		usersWithoutPassword = append(usersWithoutPassword, GetUserResponse(user))
	}

	return usersWithoutPassword
}
