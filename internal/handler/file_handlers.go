package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"project_sprint/internal/awsr"
	"project_sprint/internal/service"
	"strconv"
	"strings"
)

func UploadFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not found"})
		return
	}

	fileExtension := strings.ToLower(filepath.Ext(file.Filename))
	if fileExtension != ".jpeg" && fileExtension != ".jpg" && fileExtension != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only jpeg, jpg, png allowed."})
		return
	}

	if file.Size > 1024*100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 100KiB"})
		return
	}

	uploadPath := filepath.Join("uploads", file.Filename)

	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file locally"})
		return
	}

	fileURL, err := awsr.UploadToS3(uploadPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to S3"})
		return
	}

	storedFile, err := service.AddFile(fileURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store file in the database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":  storedFile.ID,
		"uri": storedFile.URI,
		// URL file di S3
	})
}

func GetFileHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	file, err := service.GetFileByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": file.ID, "uri": file.URI})
}

func DeleteFileHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	err = service.DeleteFile(id)
	if err != nil {
		if err.Error() == "file not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
