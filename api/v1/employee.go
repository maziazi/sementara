package v1

import (
	"github.com/gin-gonic/gin"
	"project_sprint/internal/handler"
	"project_sprint/internal/middleware"
)

func RegisterEmployeeRoutes(router *gin.RouterGroup) {

	protected := router.Group("/employee")
	protected.Use(middleware.JWTAuthMiddleware())
	router.Group("/employee").GET("/", handler.GetEmployeeHandler)
	{
		protected.POST("/", handler.CreateEmployeeHandler)
	}
}
