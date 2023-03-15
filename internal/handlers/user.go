package handlers

import (
	"errors"
	"log"
	"net/http"
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

const (
	HASH_COST = 8
)

type SignUpLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	SsoType  string `json:"sso_type"`
	Code     string `json:"code"`
}

type SignUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	User  models.UserResponse `json:"user"`
	Token string              `json:"token"`
}

type UserUpdateRequest struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Picture     string `json:"picture"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

// SignUpHandler handles the sign up request
func SignUpHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request = models.User{}

		err := c.BindJSON(&request)
		if err != nil {
			HandleError(c, http.StatusBadRequest, err)
			return
		}

		if !utils.ValidateEmail(request.Email) {
			HandleError(c, http.StatusBadRequest, errors.New("invalid email"))
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), HASH_COST)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		var user = models.User{
			Id:          id.String(),
			FirstName:   request.FirstName,
			LastName:    request.LastName,
			Username:    request.Username,
			Email:       request.Email,
			Password:    string(hashedPassword),
			PhoneNumber: request.PhoneNumber,
			Picture:     request.Picture,
			Address:     request.Address,
			IsActive:    true,
		}

		u, err := repository.InsertUser(c.Request.Context(), &user)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		// Add default role to user
		_, err = AddRoleToUser(c.Request.Context(), u.Id, "user")
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		// Generate token and save it
		token, err := utils.GenerateToken()
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		// Save token
		err = SaveVerifiedEmailToken(c.Request.Context(), u, token)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		// Send verification email
		err = SendVerificationEmail(s, u, token)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		// Send response
		signUpResponse := SignUpResponse{
			Id:    u.Id,
			Email: u.Email,
		}

		HandleSuccess(c, http.StatusCreated, "ok", signUpResponse)
	}
}

// VerifyEmailHandler handles the verification of email
func VerifiedEmailHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")

		user, err := repository.GetUserByVerifiedEmailToken(c.Request.Context(), token)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		if time.Now().After(user.VerifiedEmailTokenExpiry) {
			HandleError(c, http.StatusUnauthorized, errors.New("token expired"))
			return
		}

		user.VerifiedEmail = true

		_, err = repository.UpdateUser(c.Request.Context(), user.Id, user)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		c.Redirect(http.StatusMovedPermanently, s.Config().FrontendURL+"/auth/login")
	}
}

// ResetPasswordHandler handles the reset password request
func ResetPasswordHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request = ResetPasswordRequest{}

		err := c.BindJSON(&request)
		if err != nil {
			HandleError(c, http.StatusBadRequest, err)
			return
		}

		if !utils.ValidateEmail(request.Email) {
			HandleError(c, http.StatusBadRequest, errors.New("invalid email"))
			return
		}

		user, err := repository.GetUserByEmail(c.Request.Context(), request.Email)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		// Generate token and save it
		token, err := utils.GenerateToken()
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		// Save token
		err = SavePasswordResetToken(c.Request.Context(), user, token)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		// Send reset password email
		err = SendResetPasswordEmail(s, user, token)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		HandleSuccess(c, http.StatusOK, "ok", nil)
	}
}

// LoginHandler handles the login request
func LoginHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request = SignUpLoginRequest{}
		var userRq = models.UserResponse{}

		err := c.BindJSON(&request)
		if err != nil {
			HandleError(c, http.StatusBadRequest, err)
			return
		}

		if request.SsoType == "google" {
			// Login with google
			user, err := HandleGoogleLogin(c, s, &request)
			if err != nil {
				HandleError(c, http.StatusInternalServerError, err)
				return
			}

			userRq = models.UserResponse{
				Id:            user.Id,
				FirstName:     user.FirstName,
				LastName:      user.LastName,
				Username:      user.Username,
				Email:         user.Email,
				Picture:       user.Picture,
				Address:       user.Address,
				PhoneNumber:   user.PhoneNumber,
				Roles:         user.Roles,
				IsActive:      user.IsActive,
				VerifiedEmail: user.VerifiedEmail,
			}
		} else {
			// Login with email and password
			user, err := HandleEmailAndPasswordLogin(c, s, &request)
			if err != nil {
				HandleError(c, http.StatusInternalServerError, err)
				return
			}

			userRq = models.UserResponse{
				Id:            user.Id,
				FirstName:     user.FirstName,
				LastName:      user.LastName,
				Username:      user.Username,
				Email:         user.Email,
				Picture:       user.Picture,
				Address:       user.Address,
				PhoneNumber:   user.PhoneNumber,
				Roles:         user.Roles,
				IsActive:      user.IsActive,
				VerifiedEmail: user.VerifiedEmail,
			}
		}

		// Generate JWT token
		claims := models.AppClaims{
			UserId: userRq.Id,
			Email:  userRq.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(s.Config().JWTSecret))
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		loginResponse := LoginResponse{
			User:  userRq,
			Token: tokenString,
		}

		HandleSuccess(c, http.StatusOK, "ok", loginResponse)
	}
}

// MeHandler handles the me request
func MeHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		claims, err := DecodeToken(tokenString, s.Config().JWTSecret)
		if err != nil {
			HandleError(c, http.StatusUnauthorized, err)
			return
		}

		user, err := s.Redis().GetValue(claims.UserId)
		if err != nil {
			log.Println(err)
		}

		if user == nil {
			user, err := repository.GetUserById(c.Request.Context(), claims.UserId)
			if err != nil {
				HandleError(c, http.StatusInternalServerError, err)
				return
			}

			err = s.Redis().SetValue(user.Id, user)
			if err != nil {
				log.Println(err)
			}

			HandleSuccess(c, http.StatusOK, "ok", user)
			return
		}

		HandleSuccess(c, http.StatusOK, "ok", user)
	}
}

// UpdateUserHandler handles the update user request
func UpdateUserHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		claims, err := DecodeToken(tokenString, s.Config().JWTSecret)
		if err != nil {
			HandleError(c, http.StatusUnauthorized, err)
			return
		}

		var request = models.UserResponse{}

		err = c.BindJSON(&request)
		if err != nil {
			HandleError(c, http.StatusBadRequest, err)
			return
		}

		user, err := repository.UpdateUser(c.Request.Context(), claims.UserId, &request)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		HandleSuccess(c, http.StatusOK, "ok", user)
	}
}

// UploadPictureHandler handles the upload picture request
func UploadPictureHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		if id == "" {
			HandleError(c, http.StatusBadRequest, errors.New("invalid id"))
			return
		}

		maxSize := int64(1024 * 1024 * 5) // 5MB

		err := c.Request.ParseMultipartForm(maxSize)
		if err != nil {
			HandleError(c, http.StatusBadRequest, err)
			return
		}

		file, fileHeader, err := c.Request.FormFile("picture")
		if err != nil {
			HandleError(c, http.StatusBadRequest, err)
			return
		}

		defer file.Close()

		// reused if we're uploading many files
		fileName, err := s.S3().UploadFileToS3(file, fileHeader, id)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		fileUrl := s.S3().GenerateUrl(fileName)

		user, err := repository.GetUserById(c.Request.Context(), id)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		user.Picture = fileUrl

		u, err := repository.PartialUpdateUser(c.Request.Context(), id, user)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		HandleSuccess(c, http.StatusOK, "ok", u)
	}
}

// ListUserHandler handles the list user request
func ListUserHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := s.Redis().GetValue("users")
		if err != nil {
			log.Println(err)
		}

		if users == nil {
			users, err := repository.ListUser(c.Request.Context())
			if err != nil {
				HandleError(c, http.StatusInternalServerError, err)
				return
			}

			err = s.Redis().SetValue("users", users)
			if err != nil {
				log.Println(err)
			}

			HandleSuccess(c, http.StatusOK, "ok", users)
			return
		}

		HandleSuccess(c, http.StatusOK, "ok", users)
	}
}
