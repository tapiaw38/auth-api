package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api/internal/models"
	"github.com/tapiaw38/auth-api/internal/repository"
	"github.com/tapiaw38/auth-api/internal/server"
)

type UserRoleRequest struct {
	UserId string `json:"user_id"`
	RoleId string `json:"role_id"`
}

// InsertUserRole is the handler for the InsertUserRole endpoint
func InsertUserRole(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request = UserRoleRequest{}

		err := c.BindJSON(&request)
		if err != nil {
			HandleError(c, http.StatusBadRequest, err)
			return
		}

		var userRole = models.UserRole{
			UserId: request.UserId,
			RoleId: request.RoleId,
		}

		err = repository.InsertUserRole(c.Request.Context(), &userRole)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		HandleSuccess(c, http.StatusOK, "ok", nil)
	}
}

// DeleteUserRole is the handler for the DeleteUserRole endpoint
func DeleteUserRole(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request = UserRoleRequest{}

		err := c.BindJSON(&request)
		if err != nil {
			HandleError(c, http.StatusBadRequest, err)
			return
		}

		userRole := models.UserRole{
			UserId: request.UserId,
			RoleId: request.RoleId,
		}

		err = repository.DeleteUserRole(c.Request.Context(), &userRole)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		HandleSuccess(c, http.StatusOK, "ok", nil)
	}
}
