package v1

import (
	"github.com/gin-gonic/gin"
	"project_sprint/internal/handler"
	"project_sprint/internal/middleware"
)

func RegisterEmployeeRoutes(router *gin.RouterGroup) {

	// Rute untuk mendapatkan daftar employee tanpa autentikasi
	router.GET("/employee", handler.GetEmployeeHandler)

	// Grup rute yang memerlukan autentikasi (JWT)
	protected := router.Group("/employee")
	protected.Use(middleware.JWTAuthMiddleware())

	// Rute untuk menambah employee baru dan update data employee
	protected.POST("/", handler.CreateEmployeeHandler)
	protected.PATCH("/:identityNumber", handler.UpdateEmployeeHandler) // Menggunakan :identityNumber untuk mengupdate employee berdasarkan ID
	protected.DELETE("/:identityNumber", handler.DeleteEmployeeHandler)
}
