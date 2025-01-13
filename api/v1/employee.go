package v1

import (
	"github.com/gin-gonic/gin"
	"project_sprint/internal/handler"
	"project_sprint/internal/middleware"
)

func RegisterEmployeeRoutes(router *gin.RouterGroup) {

	protected := router.Group("/employee")
	protected.Use(middleware.JWTAuthMiddleware()) // Gunakan middleware JWT
	router.Group("/employee").GET("/", handler.GetEmployeeHandler)
	{
		//protected.GET("/", handler.)      // Ambil profil user
		// Update profil user
		protected.POST("/", handler.CreateEmployeeHandler)
	}
}
