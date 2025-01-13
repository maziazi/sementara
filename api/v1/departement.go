package v1

import (
	"github.com/gin-gonic/gin"
	"project_sprint/internal/handler"
)

// RegisterDepartmentRoutes mendaftarkan endpoint Department
func RegisterDepartmentRoutes(router *gin.RouterGroup) {
	departmentGroup := router.Group("/departments")
	{
		departmentGroup.POST("/", handler.CreateDepartmentHandler)
		departmentGroup.GET("/", handler.GetDepartments)
	}
}

//func RegisterDepartmentRoutes(router *mux.Router) {
//	router.HandleFunc("/departements", handler.CreateDepartmentHandler).Methods(http.MethodPost)
//	router.HandleFunc("/departements", handler.GetDepartments).Methods(http.MethodGet)
//}
