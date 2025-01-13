package v1

import (
	"github.com/gin-gonic/gin"
	"project_sprint/internal/handler"
)

func RegisterFileRoutes(router *gin.RouterGroup) {
	//TODO FILE BELUM DIKASIH PROTECTION
	router.POST("/files", handler.UploadFileHandler)
	router.GET("/files/:id", handler.GetFileHandler)
	router.DELETE("/files/:id", handler.DeleteFileHandler)

}
