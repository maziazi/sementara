package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"project_sprint/internal/service"
)

// Struct untuk request pembuatan employee
type EmployeeRequest struct {
	IdentityNumber   string `json:"identityNumber" binding:"required,min=5,max=33"`
	Name             string `json:"name" binding:"required,min=4,max=33"`
	EmployeeImageUri string `json:"employeeImageUri" binding:"required,url"`
	Gender           string `json:"gender" binding:"required,oneof=male female"`
	DepartmentId     string `json:"departmentId" binding:"required"`
}

// Handler untuk membuat employee baru
func CreateEmployeeHandler(c *gin.Context) {
	var req EmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Konversi departmentId dari string ke int
	departmentId, err := strconv.Atoi(req.DepartmentId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid departmentId: must be an integer"})
		return
	}

	// Panggil service untuk membuat employee
	employee, err := service.CreateEmployee(req.IdentityNumber, req.Name, req.EmployeeImageUri, req.Gender, departmentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"identityNumber":   employee.IdentityNumber,
		"name":             employee.Name,
		"employeeImageUri": employee.EmployeeImageURI,
		"gender":           employee.Gender,
		"departmentId":     strconv.Itoa(employee.DepartmentID), // Konversi ke string
	})
}

// Handler untuk mendapatkan daftar employees
func GetEmployeeHandler(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "5")
	offsetParam := c.DefaultQuery("offset", "0")
	identityNumber := c.DefaultQuery("identityNumber", "")
	name := c.DefaultQuery("name", "")
	gender := c.DefaultQuery("gender", "")
	departmentIdParam := c.DefaultQuery("departmentId", "0")

	// Konversi limit dan offset ke integer
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 5
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Validasi gender (hanya boleh "male" atau "female")
	var validGender string
	if gender != "" && gender != "male" && gender != "female" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gender"})
		return
	}
	validGender = gender

	// Konversi departmentId ke integer (jika valid)
	departmentId, err := strconv.Atoi(departmentIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid departmentId: must be an integer"})
		return
	}

	// Panggil service untuk mengambil data employees
	employees, err := service.GetEmployees(limit, offset, identityNumber, name, validGender, departmentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ubah `departmentId` dari int ke string dalam respons JSON
	var response []map[string]interface{}
	for _, emp := range employees {
		response = append(response, map[string]interface{}{
			"identityNumber":   emp.IdentityNumber,
			"name":             emp.Name,
			"employeeImageUri": emp.EmployeeImageURI,
			"gender":           emp.Gender,
			"departmentId":     strconv.Itoa(emp.DepartmentID), // Konversi ke string
		})
	}

	c.JSON(http.StatusOK, response)
}

// Struct untuk request update employee
type EmployeeUpdateRequest struct {
	Name             *string `json:"name"`
	EmployeeImageUri *string `json:"employeeImageUri"`
	Gender           *string `json:"gender"`
	DepartmentId     *string `json:"departmentId"`
}

// Handler untuk update employee
func UpdateEmployeeHandler(c *gin.Context) {
	identityNumber := c.Param("identityNumber")
	if identityNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "identityNumber is required"})
		return
	}

	var req EmployeeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Konversi departmentId jika ada
	var departmentId *int
	if req.DepartmentId != nil {
		id, err := strconv.Atoi(*req.DepartmentId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid departmentId: must be an integer"})
			return
		}
		departmentId = &id
	}

	// Panggil service untuk update employee
	updatedEmployee, err := service.UpdateEmployee(
		identityNumber,
		// Periksa jika nil atau tidak untuk setiap field
		nilString(req.Name),
		nilString(req.EmployeeImageUri),
		nilString(req.Gender),
		departmentId,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"identityNumber":   updatedEmployee.IdentityNumber,
		"name":             updatedEmployee.Name,
		"employeeImageUri": updatedEmployee.EmployeeImageURI,
		"gender":           updatedEmployee.Gender,
		"departmentId":     strconv.Itoa(updatedEmployee.DepartmentID), // Konversi ke string
	})
}

// Fungsi untuk menangani nil pointer pada string
func nilString(s *string) string {
	if s == nil {
		return "" // kembalikan string kosong jika nil
	}
	return *s // kembalikan nilai string jika tidak nil
}

// Handler untuk menghapus employee berdasarkan identityNumber
func DeleteEmployeeHandler(c *gin.Context) {
	identityNumber := c.Param("identityNumber")
	if identityNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "identityNumber is required"})
		return
	}

	// Panggil service untuk menghapus employee
	err := service.DeleteEmployee(identityNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kembalikan response sukses
	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}
