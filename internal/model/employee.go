package model

type Employee struct {
	ID               int    `json:"id"`
	IdentityNumber   string `json:"identityNumber"`
	Name             string `json:"name"`
	EmployeeImageURI string `json:"employeeImageUri"`
	Gender           string `json:"gender"`
	DepartmentID     int    `json:"departmentId"` // Ubah ke int
}
