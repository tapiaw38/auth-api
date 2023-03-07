package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"github.com/tapiaw38/auth-api/internal/models"
	"github.com/tapiaw38/auth-api/internal/repository"
	"github.com/tapiaw38/auth-api/internal/server"
)

type RoleRequest struct {
	Name string `json:"name"`
}

// InsertRoleHandler handles the insert role request
func InsertRoleHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request = RoleRequest{}

		err := c.BindJSON(&request)
		if err != nil {
			HandleError(c, http.StatusBadRequest, err)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		var role = models.Role{
			Id:   id.String(),
			Name: request.Name,
		}

		rl, err := repository.InsertRole(c.Request.Context(), &role)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		HandleSuccess(c, http.StatusOK, "ok", rl)
	}
}

// GetRoleByNameHandler handles the get role by name request
func GetRoleByNameHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			HandleError(c, http.StatusBadRequest, errors.New("name is required"))
			return
		}

		role, err := repository.GetRoleByName(c.Request.Context(), name)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		HandleSuccess(c, http.StatusOK, "ok", role)
	}
}

// GetRoleByIdHandler handles the get role by id request
func GetRoleByIdHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			HandleError(c, http.StatusBadRequest, errors.New("id is required"))
			return
		}

		role, err := repository.GetRoleById(c.Request.Context(), id)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		HandleSuccess(c, http.StatusOK, "ok", role)
	}
}

// UpdateRoleHandler handles the update role request
func UpdateRoleHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request = RoleRequest{}

		err := c.BindJSON(&request)
		if err != nil {
			HandleError(c, http.StatusBadRequest, err)
			return
		}

		id := c.Param("id")
		if id == "" {
			HandleError(c, http.StatusBadRequest, errors.New("id is required"))
			return
		}

		var role = models.Role{
			Id:   id,
			Name: request.Name,
		}

		rl, err := repository.UpdateRole(c.Request.Context(), &role)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		HandleSuccess(c, http.StatusOK, "ok", rl)
	}
}

// DeleteRoleHandler handles the delete role request
func DeleteRoleHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			HandleError(c, http.StatusBadRequest, errors.New("id is required"))
			return
		}

		role, err := repository.DeleteRole(c.Request.Context(), id)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		HandleSuccess(c, http.StatusOK, "ok", role)
	}
}

// ListRolesHandler handles the list roles request
func ListRoleHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, err := repository.ListRole(c.Request.Context())
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}

		HandleSuccess(c, http.StatusOK, "ok", roles)
	}
}
