package v1

import (
	"github.com/gin-gonic/gin"
	"project_sprint/internal/handler"
	"project_sprint/internal/middleware"
)

func RegisterDepartmentRoutes(router *gin.RouterGroup) {

	/*TODO(
		AUTH UNTUK DEPARTMENT
		POST
	GET
	PATCH
	DELETE)

	*/
	router.Group("/department").GET("/", handler.GetDepartments)
	protected := router.Group("/department")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.POST("/", handler.CreateDepartmentHandler)
		protected.PATCH("/:id", handler.PatchDepartment)
	}
}
