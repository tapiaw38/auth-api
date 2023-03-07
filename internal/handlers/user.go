package handlers

import (
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

// SignUpHandler handles the sign up request
func SignUpHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request = models.User{}

		err := c.BindJSON(&request)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusBadRequest, response)
			return
		}

		if !utils.ValidateEmail(request.Email) {
			response := NewResponse(Error, "Invalid email", nil)
			ResponseWithJson(c, http.StatusBadRequest, response)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), HASH_COST)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
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
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		// Add default role to user
		_, err = AddRoleToUser(c.Request.Context(), u.Id, "user")
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		// Generate token and save it
		token, err := GenerateToken(u)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		// Send email verification
		channel := make(chan error)
		subjet := "Bienvenido a Mi Tur"

		variables := map[string]interface{}{
			"name": u.FirstName + " " + u.LastName,
			"link": s.Config().Host + "/auth/verify-email?token=" + token,
		}

		go s.Mail().SendEmail(u.Email, subjet, "email_verification", variables, channel)

		err = <-channel
		if err != nil {
			log.Println(err)
		}

		// Send response
		signUpResponse := SignUpResponse{
			Id:    u.Id,
			Email: u.Email,
		}

		response := NewResponse(Message, "ok", signUpResponse)
		ResponseWithJson(c, http.StatusCreated, response)
	}
}

func VerifiedEmailHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")

		user, err := repository.GetUserByToken(c.Request.Context(), token)
		if err != nil {
			log.Println("Error getting user by token: ", err)
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusBadRequest, response)
			return
		}

		if time.Now().After(user.TokenExpiry) {
			response := NewResponse(Error, "Token expired", nil)
			ResponseWithJson(c, http.StatusBadRequest, response)
			return
		}

		user.VerifiedEmail = true

		_, err = repository.UpdateUser(c.Request.Context(), user.Id, user)
		if err != nil {
			log.Println("Error updating user: ", err)
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		c.Redirect(http.StatusMovedPermanently, s.Config().FrontendURL+"/auth/login")
	}
}

// LoginHandler handles the login request
func LoginHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request = SignUpLoginRequest{}
		var userRq = models.UserResponse{}

		err := c.BindJSON(&request)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		if request.SsoType == "google" {
			// Login with google
			token, err := s.Google().ExchangeCode(c.Request.Context(), request.Code)
			if err != nil {
				response := NewResponse(Error, err.Error(), nil)
				ResponseWithJson(c, http.StatusBadRequest, response)
				return
			}

			userInfo, err := s.Google().GetUserInfo(c.Request.Context(), token)
			if err != nil {
				response := NewResponse(Error, err.Error(), nil)
				ResponseWithJson(c, http.StatusBadRequest, response)
				return
			}

			user, err := repository.GetUserByEmailSocial(c.Request.Context(), userInfo.Email)
			if err != nil || user == nil {
				// If user not registered, register user
				id, err := ksuid.NewRandom()
				if err != nil {
					response := NewResponse(Error, err.Error(), nil)
					ResponseWithJson(c, http.StatusInternalServerError, response)
					return
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
					response := NewResponse(Error, err.Error(), nil)
					ResponseWithJson(c, http.StatusInternalServerError, response)
					return
				}

				user, err := AddRoleToUser(c.Request.Context(), user.Id, "user")
				if err != nil {
					response := NewResponse(Error, err.Error(), nil)
					ResponseWithJson(c, http.StatusInternalServerError, response)
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
					IsActive:      user.IsActive,
					Roles:         user.Roles,
					VerifiedEmail: user.VerifiedEmail,
				}

			} else {
				// If user already registered, update user info
				if user.Picture == "" || !user.VerifiedEmail {
					userUpdate := models.UserResponse{
						Picture:       userInfo.Picture,
						IsActive:      user.IsActive,
						VerifiedEmail: userInfo.VerifiedEmail,
					}

					user, err = repository.PartialUpdateUser(c.Request.Context(), user.Id, &userUpdate)
					if err != nil {
						response := NewResponse(Error, err.Error(), nil)
						ResponseWithJson(c, http.StatusInternalServerError, response)
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
			}
		} else {
			// Login with email and password
			user, err := repository.GetUserByEmail(c.Request.Context(), request.Email)
			if err != nil {
				response := NewResponse(Error, err.Error(), nil)
				ResponseWithJson(c, http.StatusInternalServerError, response)
				return
			}

			if user == nil {
				response := NewResponse(Error, "Invalid Credentials", nil)
				ResponseWithJson(c, http.StatusUnauthorized, response)
				return
			}

			if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
				response := NewResponse(Error, "Invalid Credentials", nil)
				ResponseWithJson(c, http.StatusUnauthorized, response)
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
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		loginResponse := LoginResponse{
			User:  userRq,
			Token: tokenString,
		}

		response := NewResponse(Message, "ok", loginResponse)
		ResponseWithJson(c, http.StatusOK, response)
	}
}

// MeHandler handles the me request
func MeHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusUnauthorized, response)
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {

			user, err := s.Redis().GetValue(claims.UserId)
			if err != nil {
				log.Println(err)
			}

			if user == nil {
				user, err := repository.GetUserById(c.Request.Context(), claims.UserId)
				if err != nil {
					response := NewResponse(Error, err.Error(), nil)
					ResponseWithJson(c, http.StatusInternalServerError, response)
					return
				}

				err = s.Redis().SetValue(user.Id, user)
				if err != nil {
					log.Println(err)
				}

				response := NewResponse(Message, "ok", user)
				ResponseWithJson(c, http.StatusCreated, response)
				return
			}
			response := NewResponse(Message, "ok", user)
			ResponseWithJson(c, http.StatusCreated, response)
			return
		}

		response := NewResponse(Error, "Invalid Token", nil)
		ResponseWithJson(c, http.StatusInternalServerError, response)
	}
}

// UpdateUserHandler handles the update user request
func UpdateUserHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusUnauthorized, response)
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			var request = models.UserResponse{}

			err := c.BindJSON(&request)
			if err != nil {
				response := NewResponse(Error, err.Error(), nil)
				ResponseWithJson(c, http.StatusInternalServerError, response)
				return
			}

			user, err := repository.UpdateUser(c.Request.Context(), claims.UserId, &request)
			if err != nil {
				response := NewResponse(Error, err.Error(), nil)
				ResponseWithJson(c, http.StatusInternalServerError, response)
				return
			}

			response := NewResponse(Message, "ok", user)
			ResponseWithJson(c, http.StatusCreated, response)
		}
		response := NewResponse(Error, "Invalid Token", nil)
		ResponseWithJson(c, http.StatusInternalServerError, response)
	}
}

// UploadPictureHandler handles the upload picture request
func UploadPictureHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		if id == "" {
			response := NewResponse(Error, "Invalid id", nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		maxSize := int64(1024 * 1024 * 5) // 5MB

		err := c.Request.ParseMultipartForm(maxSize)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		file, fileHeader, err := c.Request.FormFile("picture")
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		defer file.Close()

		// reused if we're uploading many files
		fileName, err := s.S3().UploadFileToS3(file, fileHeader, id)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		fileUrl := s.S3().GenerateUrl(fileName)

		user, err := repository.GetUserById(c.Request.Context(), id)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		user.Picture = fileUrl

		u, err := repository.PartialUpdateUser(c.Request.Context(), id, user)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		response := NewResponse(Message, "ok", u)
		ResponseWithJson(c, http.StatusOK, response)
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
				response := NewResponse(Error, err.Error(), nil)
				ResponseWithJson(c, http.StatusInternalServerError, response)
				return
			}

			err = s.Redis().SetValue("users", users)
			if err != nil {
				log.Println(err)
			}

			response := NewResponse(Message, "ok", users)
			ResponseWithJson(c, http.StatusCreated, response)
			return
		}

		response := NewResponse(Message, "ok", users)
		ResponseWithJson(c, http.StatusCreated, response)
	}
}
