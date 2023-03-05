package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"github.com/tapiaw38/auth-api/models"
	"github.com/tapiaw38/auth-api/repository"
	"github.com/tapiaw38/auth-api/server"
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
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusBadRequest, response)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		var role = models.Role{
			Id:   id.String(),
			Name: request.Name,
		}

		rl, err := repository.InsertRole(c.Request.Context(), &role)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		response := NewResponse(Message, "ok", rl)
		ResponseWithJson(c, http.StatusCreated, response)
	}
}

// GetRoleByNameHandler handles the get role by name request
func GetRoleByNameHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			response := NewResponse(Error, "name is required", nil)
			ResponseWithJson(c, http.StatusBadRequest, response)
			return
		}

		role, err := repository.GetRoleByName(c.Request.Context(), name)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		response := NewResponse(Message, "ok", role)
		ResponseWithJson(c, http.StatusOK, response)
	}
}

// GetRoleByIdHandler handles the get role by id request
func GetRoleByIdHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			response := NewResponse(Error, "id is required", nil)
			ResponseWithJson(c, http.StatusBadRequest, response)
			return
		}

		role, err := repository.GetRoleById(c.Request.Context(), id)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		response := NewResponse(Message, "ok", role)
		ResponseWithJson(c, http.StatusOK, response)
	}
}

// UpdateRoleHandler handles the update role request
func UpdateRoleHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request = RoleRequest{}

		err := c.BindJSON(&request)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusBadRequest, response)
			return
		}

		id := c.Param("id")
		if id == "" {
			response := NewResponse(Error, "id is required", nil)
			ResponseWithJson(c, http.StatusBadRequest, response)
			return
		}

		var role = models.Role{
			Id:   id,
			Name: request.Name,
		}

		rl, err := repository.UpdateRole(c.Request.Context(), &role)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		response := NewResponse(Message, "ok", rl)
		ResponseWithJson(c, http.StatusOK, response)
	}
}

// DeleteRoleHandler handles the delete role request
func DeleteRoleHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			response := NewResponse(Error, "id is required", nil)
			ResponseWithJson(c, http.StatusBadRequest, response)
			return
		}

		role, err := repository.DeleteRole(c.Request.Context(), id)
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		response := NewResponse(Message, "ok", role)
		ResponseWithJson(c, http.StatusOK, response)
	}
}

// ListRolesHandler handles the list roles request
func ListRoleHandler(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, err := repository.ListRole(c.Request.Context())
		if err != nil {
			response := NewResponse(Error, err.Error(), nil)
			ResponseWithJson(c, http.StatusInternalServerError, response)
			return
		}

		response := NewResponse(Message, "ok", roles)
		ResponseWithJson(c, http.StatusOK, response)
	}
}
