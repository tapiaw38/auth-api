package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api/models"
	"github.com/tapiaw38/auth-api/repository"
	"github.com/tapiaw38/auth-api/server"
)

type UserRoleRequest struct {
	UserId string `json:"user_id"`
	RoleId string `json:"role_id"`
}

// InsertUserRole is the handler for the InsertUserRole endpoint
func InsertUserRole(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request = UserRoleRequest{}

		log.Println("InsertUserRole", request)

		err := c.BindJSON(&request)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusBadRequest, response)
			return
		}

		var userRole = models.UserRole{
			UserId: request.UserId,
			RoleId: request.RoleId,
		}

		err = repository.InsertUserRole(c.Request.Context(), &userRole)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		response := NewResponse(Message, "ok", nil)
		ResponseWithJson(c, http.StatusCreated, response)
	}
}

// DeleteUserRole is the handler for the DeleteUserRole endpoint
func DeleteUserRole(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request = UserRoleRequest{}

		err := c.BindJSON(&request)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusBadRequest, response)
			return
		}

		userRole := models.UserRole{
			UserId: request.UserId,
			RoleId: request.RoleId,
		}

		err = repository.DeleteUserRole(c.Request.Context(), &userRole)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		response := NewResponse(Message, "ok", nil)
		ResponseWithJson(c, http.StatusOK, response)
	}
}
