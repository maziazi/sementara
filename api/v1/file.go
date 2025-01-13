package v1

import (
	"github.com/gin-gonic/gin"
	"project_sprint/internal/handler"
)

func RegisterFileRoutes(router *gin.RouterGroup) {

	router.POST("/files", handler.UploadFileHandler)
	router.GET("/files/:id", handler.GetFileHandler)
	router.DELETE("/files/:id", handler.DeleteFileHandler)

}
