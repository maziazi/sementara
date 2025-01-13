package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project_sprint/internal/service"
	"strconv"
)

func CreateEmployeeHandler(c *gin.Context) {
	type EmployeeRequest struct {
		IdentityNumber   string `json:"identityNumber" binding:"required,min=5,max=33"`
		Name             string `json:"name" binding:"required,min=4,max=33"`
		EmployeeImageUri string `json:"employeeImageUri" binding:"required,url"`
		Gender           string `json:"gender" binding:"required,oneof=male female"`
		DepartmentId     string `json:"departmentId" binding:"required"`
	}

	var req EmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//if !service.IsDepartmentValid(req.DepartmentId) {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid departmentId"})
	//	return
	//}

	employee, err := service.CreateEmployee(req.IdentityNumber, req.Name, req.EmployeeImageUri, req.Gender, req.DepartmentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":               employee.ID,
		"identityNumber":   employee.IdentityNumber,
		"name":             employee.Name,
		"employeeImageUri": employee.EmployeeImageURI,
		"gender":           employee.Gender,
		"departmentId":     employee.DepartmentID,
	})
}

func GetEmployeeHandler(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "5")
	offsetParam := c.DefaultQuery("offset", "0")
	identityNumber := c.DefaultQuery("identityNumber", "")
	name := c.DefaultQuery("name", "")
	gender := c.DefaultQuery("gender", "")
	departmentId := c.DefaultQuery("departmentId", "0")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 5
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	var validGender string
	if gender != "" {
		if gender == "male" || gender == "female" {
			validGender = gender
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gender"})
			return
		}
	}

	employees, err := service.GetEmployees(limit, offset, identityNumber, name, validGender, departmentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employees)
}
