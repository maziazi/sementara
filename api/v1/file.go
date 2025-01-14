package v1

import (
	"github.com/gin-gonic/gin"
	"project_sprint/internal/handler"
	"project_sprint/internal/middleware"
)

func RegisterFileRoutes(router *gin.RouterGroup) {
	//TODO FILE BELUM DIKASIH PROTECTION
	router.Use(middleware.JWTAuthMiddleware()).POST("/file", handler.UploadFileHandler)
	router.GET("/file/:id", handler.GetFileHandler)
	router.DELETE("/file/:id", handler.DeleteFileHandler)

}
