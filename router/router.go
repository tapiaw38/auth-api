package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api/handlers"
	"github.com/tapiaw38/auth-api/middleware"
	"github.com/tapiaw38/auth-api/server"
)

// BinderRoutes mounts the routes and the middleware
func BinderRoutes(s server.Server, router *gin.Engine) {

	authRoute := router.Group("/auth/")
	authRoute.POST("signup", handlers.SignUpHandler(s))
	authRoute.POST("login", handlers.LoginHandler(s))

	// mount the middleware
	router.Use(middleware.CheckAuthMiddleware(s))

	// User routes
	userRoute := router.Group("/users/")
	userRoute.GET("me", handlers.MeHandler(s))
	userRoute.PUT(":id", handlers.UpdateUserHandler(s))
	userRoute.PUT("picture/:id", handlers.UploadPictureHandler(s))
	userRoute.GET("list", handlers.ListUserHandler(s))

	// Role routes
	roleRoute := router.Group("/roles/")
	roleRoute.POST("new", handlers.InsertRoleHandler(s))
	roleRoute.GET("list", handlers.ListRoleHandler(s))
	roleRoute.GET(":id", handlers.GetRoleByIdHandler(s))
	roleRoute.PUT(":id", handlers.UpdateRoleHandler(s))
	roleRoute.DELETE(":id", handlers.DeleteRoleHandler(s))

	// User Role routes
	userRoleRoute := router.Group("/user_roles/")
	userRoleRoute.POST("new", handlers.InsertUserRole(s))
	userRoleRoute.DELETE("delete", handlers.DeleteUserRole(s))
}
