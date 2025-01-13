package v1

import (
	"github.com/gin-gonic/gin"
	"project_sprint/internal/handler"
	"project_sprint/internal/middleware"
)

func RegisterManagerRoutes(router *gin.RouterGroup) {

	protected := router.Group("/user")
	protected.Use(middleware.JWTAuthMiddleware()) // Gunakan middleware JWT
	{
		protected.GET("/", handler.GetUserProfileHandler)      // Ambil profil user
		protected.PATCH("/", handler.UpdateUserProfileHandler) // Update profil user
	}
}
