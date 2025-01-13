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

// UploadFileHandler untuk menyimpan URI file ke database
//
//	func UploadFileHandler(c *gin.Context) {
//		// Mendapatkan file dari form-data
//		file, err := c.FormFile("file") // 'file' adalah key di form-data
//		if err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
//			return
//		}
//
//		// Validasi jenis file
//		fileExtension := strings.ToLower(filepath.Ext(file.Filename))
//		if fileExtension != ".jpeg" && fileExtension != ".jpg" && fileExtension != ".png" {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only jpeg, jpg, png allowed."})
//			return
//		}
//
//		// Validasi ukuran file (maksimal 100 KiB)
//		if file.Size > 1024*100 {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 100KiB"})
//			return
//		}
//
//		// Tentukan lokasi penyimpanan sementara untuk file
//		dst := filepath.Join("uploads", file.Filename)
//		if err := c.SaveUploadedFile(file, dst); err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
//			return
//		}
//
//		// Simpan file ke layanan cloud storage (misalnya, S3) jika diperlukan
//		// Misalnya, simpan URI file ke dalam database setelah di-upload ke S3
//		// Ganti dengan logic untuk menyimpan file ke S3, dan dapatkan URI-nya.
//		uri := fmt.Sprintf("/uploads/%s", file.Filename) // Contoh URI
//
//		// Simpan URI file ke database
//		storedFile, err := service.AddFile(uri)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store file in the database"})
//			return
//		}
//
//		// Kembalikan response dengan informasi file
//		c.JSON(http.StatusCreated, gin.H{
//			"id":  storedFile.ID,
//			"uri": storedFile.URI,
//		})
//	}
func UploadFileHandler(c *gin.Context) {
	// Mengambil file dari form-data
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not found"})
		return
	}

	// Validasi jenis file
	fileExtension := strings.ToLower(filepath.Ext(file.Filename))
	if fileExtension != ".jpeg" && fileExtension != ".jpg" && fileExtension != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only jpeg, jpg, png allowed."})
		return
	}

	// Validasi ukuran file (maksimal 100 KiB)
	if file.Size > 1024*100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 100KiB"})
		return
	}

	// Tentukan path untuk menyimpan file sementara
	uploadPath := filepath.Join("uploads", file.Filename)

	// Menyimpan file sementara di server lokal sebelum di-upload ke S3
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file locally"})
		return
	}

	// Upload file ke S3
	fileURL, err := awsr.UploadToS3(uploadPath) // Fungsi untuk upload ke S3
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to S3"})
		return
	}

	// Simpan URI file yang sudah di-upload ke S3 ke dalam database
	storedFile, err := service.AddFile(fileURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store file in the database"})
		return
	}

	// Kembalikan response dengan informasi file
	c.JSON(http.StatusCreated, gin.H{
		"id":  storedFile.ID,
		"uri": storedFile.URI,
		// URL file di S3
	})
}

// GetFileHandler untuk mengambil file berdasarkan ID
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

// DeleteFileHandler untuk menghapus file berdasarkan ID
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
