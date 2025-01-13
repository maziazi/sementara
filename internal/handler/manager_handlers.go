package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"project_sprint/internal/middleware"
	"project_sprint/internal/service"
)

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
	Action   string `json:"action" binding:"required,oneof=create login"`
}

func AuthHandler(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Action == "create" {
		user, err := service.RegisterUser(req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}

		token, _ := middleware.GenerateToken(user.Email, user.ID)
		c.JSON(http.StatusCreated, gin.H{"email": user.Email, "token": token})
		return
	}

	if req.Action == "login" {
		user, err := service.AuthenticateManager(req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid email or password"})
			return
		}

		token, _ := middleware.GenerateToken(user.Email, user.ID)
		c.JSON(http.StatusOK, gin.H{"email": user.Email, "userID": user.ID, "token": token})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
}

func GetUserProfileHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDInt, ok := userID.(int)
	fmt.Println(ok)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := service.GetUserProfile(userIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email":           user.Email,
		"name":            user.Name,
		"userImageUri":    user.ManagerImageURI,
		"companyName":     user.CompanyName,
		"companyImageUri": user.CompanyImageURI,
	})

}

func UpdateUserProfileHandler(c *gin.Context) {
	var req struct {
		Email           string `json:"email" binding:"omitempty,email"`
		Name            string `json:"name" binding:"omitempty,min=4,max=52"`
		UserImageUri    string `json:"userImageUri" binding:"omitempty,uri"`
		CompanyName     string `json:"companyName" binding:"omitempty,min=4,max=52"`
		CompanyImageUri string `json:"companyImageUri" binding:"omitempty,uri"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDInt, ok := userID.(int)
	fmt.Println(ok)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := service.UpdateUserProfile(userIDInt, req.Email, req.Name, req.UserImageUri, req.CompanyName, req.CompanyImageUri)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email":           user.Email,
		"name":            user.Name,
		"userImageUri":    user.ManagerImageURI,
		"companyName":     user.CompanyName,
		"companyImageUri": user.CompanyImageURI,
	})
}
