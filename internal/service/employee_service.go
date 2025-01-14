package service

import (
	"context"
	"fmt"
	_ "strconv"
	"strings"

	"project_sprint/internal/model"
	"project_sprint/pkg/database"
)

// Fungsi untuk membuat employee baru
func CreateEmployee(identityNumber, name, employeeImageUri, gender string, departmentId int) (*model.Employee, error) {
	_, err := database.GetDBPool().
		Exec(context.Background(),
			"INSERT INTO employees (identitynumber, name, employee_image_uri, gender, departement_id) VALUES ($1, $2, $3, $4, $5)",
			identityNumber, name, employeeImageUri, gender, departmentId)
	if err != nil {
		return nil, fmt.Errorf("failed to create employee: %v", err)
	}

	return &model.Employee{
		IdentityNumber:   identityNumber,
		Name:             name,
		EmployeeImageURI: employeeImageUri,
		Gender:           gender,
		DepartmentID:     departmentId,
	}, nil
}

// Fungsi untuk mendapatkan daftar employees dengan filter
func GetEmployees(limit, offset int, identityNumber, name, gender string, departmentId int) ([]model.Employee, error) {

	db := database.GetDBPool()

	query := "SELECT id, identitynumber, name, employee_image_uri, gender, departement_id FROM employees WHERE 1=1"
	var args []interface{}
	argCount := 1

	// Filter berdasarkan identityNumber (prefix matching)
	if identityNumber != "" {
		query += fmt.Sprintf(" AND identitynumber ILIKE $%d", argCount)
		args = append(args, identityNumber+"%")
		argCount++
	}

	// Filter berdasarkan name (contains matching)
	if name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argCount)
		args = append(args, "%"+name+"%")
		argCount++
	}

	// Filter berdasarkan gender
	if gender != "" {
		query += fmt.Sprintf(" AND gender = $%d", argCount)
		args = append(args, gender)
		argCount++
	}

	// Filter berdasarkan departmentId
	if departmentId != 0 {
		query += fmt.Sprintf(" AND department_id = $%d", argCount)
		args = append(args, departmentId)
		argCount++
	}

	// Menambahkan limit dan offset
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	// Eksekusi query
	rows, err := db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees: %v", err)
	}
	defer rows.Close()

	var employees []model.Employee
	for rows.Next() {
		var emp model.Employee
		if err := rows.Scan(&emp.ID, &emp.IdentityNumber, &emp.Name, &emp.EmployeeImageURI, &emp.Gender, &emp.DepartmentID); err != nil {
			return nil, fmt.Errorf("failed to scan employee: %v", err)
		}
		employees = append(employees, emp)
	}

	return employees, nil
}

// Fungsi untuk mengupdate employee berdasarkan identityNumber
func UpdateEmployee(identityNumber, name, employeeImageUri, gender string, departmentId *int) (*model.Employee, error) {
	// Memulai query update dasar
	query := "UPDATE employees SET "
	var args []interface{}
	argCount := 1

	// Tambahkan field yang ingin diupdate jika nilainya ada
	if name != "" {
		query += fmt.Sprintf("name = $%d, ", argCount)
		args = append(args, name)
		argCount++
	}
	if employeeImageUri != "" {
		query += fmt.Sprintf("employee_image_uri = $%d, ", argCount)
		args = append(args, employeeImageUri)
		argCount++
	}
	if gender != "" {
		query += fmt.Sprintf("gender = $%d, ", argCount)
		args = append(args, gender)
		argCount++
	}
	if departmentId != nil {
		query += fmt.Sprintf("departement_id = $%d, ", argCount)
		args = append(args, *departmentId)
		argCount++
	}

	// Hapus koma terakhir jika ada (untuk menghindari error)
	query = strings.TrimSuffix(query, ", ")

	// Tambahkan kondisi WHERE untuk identityNumber
	query += fmt.Sprintf(" WHERE identitynumber = $%d RETURNING id, identitynumber, name, employee_image_uri, gender, departement_id", argCount)
	args = append(args, identityNumber)

	// Eksekusi query
	db := database.GetDBPool()
	row := db.QueryRow(context.Background(), query, args...)

	var updatedEmployee model.Employee
	if err := row.Scan(&updatedEmployee.ID, &updatedEmployee.IdentityNumber, &updatedEmployee.Name, &updatedEmployee.EmployeeImageURI, &updatedEmployee.Gender, &updatedEmployee.DepartmentID); err != nil {
		return nil, fmt.Errorf("failed to update employee: %v", err)
	}

	return &updatedEmployee, nil
}

// Fungsi untuk menghapus employee berdasarkan identityNumber
func DeleteEmployee(identityNumber string) error {
	db := database.GetDBPool()

	// Eksekusi query DELETE berdasarkan identityNumber
	_, err := db.Exec(context.Background(),
		"DELETE FROM employees WHERE identitynumber = $1", identityNumber)
	if err != nil {
		return fmt.Errorf("failed to delete employee: %v", err)
	}

	return nil
}
