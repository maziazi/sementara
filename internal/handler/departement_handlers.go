package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project_sprint/internal/model"
	"project_sprint/internal/service"
	"project_sprint/internal/utils"
	"strconv"
)

func CreateDepartmentHandler(c *gin.Context) {
	var dept model.Department

	if err := c.ShouldBindJSON(&dept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := utils.ValidateDepartment(dept.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error: " + err.Error()})
		return
	}

	department, err := service.CreateDepartment(dept)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create department" + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, department)
}

func GetDepartments(c *gin.Context) {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	nameFilter := c.Query("name")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Menambahkan wildcard % untuk prefix dan suffix matching
	if nameFilter != "" {
		nameFilter = "%" + nameFilter + "%"
	}

	departments, err := service.GetDepartments(limit, offset, nameFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, departments)
}

func PatchDepartment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id: " + idParam})
		return
	}

	var dept model.Department
	if err := c.ShouldBindJSON(&dept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validasi nama department
	err = utils.ValidateDepartment(dept.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error: " + err.Error()})
		return
	}

	// Panggil service untuk update department
	updatedDepartment, err := service.PatchDepartment(id, dept)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update department: " + err.Error()})
		return
	}

	// Response sukses
	c.JSON(http.StatusOK, updatedDepartment)
}
