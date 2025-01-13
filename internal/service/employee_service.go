package service

import (
	"context"
	"fmt"
	"project_sprint/internal/model"
	"project_sprint/pkg/database"
	_ "project_sprint/pkg/database"
	"strconv"
)

// Fungsi service untuk membuat employee baru
func CreateEmployee(identityNumber, name, employeeImageUri, gender, departmentId string) (*model.Employee, error) {
	// Cek apakah departmentId valid
	var str, _ = strconv.Atoi(departmentId)
	if !IsDepartmentValid(str) {
		return nil, fmt.Errorf("invalid departmentId")
	}

	// Simpan employee ke database
	// Implementasi menyimpan employee (ini contoh, disesuaikan dengan database Anda)
	// Biasanya menggunakan db.Exec atau db.QueryRow untuk insert dan mendapatkan ID baru
	// Misalnya, kita bisa membuat query seperti ini:
	_, err := database.GetDBPool().
		Exec(context.Background(), "INSERT INTO employees (identitynumber, name, employee_image_uri, gender, departement_id) VALUES ($1, $2, $3, $4, $5)",
			identityNumber, name, employeeImageUri, gender, departmentId)
	if err != nil {
		return nil, fmt.Errorf("failed to create employee: %v", err)
	}

	// Kembalikan employee yang baru saja dibuat (misalnya ID didapat dari insert)
	// Anda bisa mengambil ID baru setelah insert untuk mengisi ID yang benar
	return &model.Employee{
		ID:               0, // ID baru akan didapatkan setelah INSERT
		IdentityNumber:   identityNumber,
		Name:             name,
		EmployeeImageURI: employeeImageUri,
		Gender:           gender,
		DepartmentID:     departmentId,
	}, nil
}

func GetEmployees(limit, offset int, identityNumber, name, gender, departmentId string) ([]model.Employee, error) {
	db := database.GetDBPool()

	// Membuat query dasar untuk mencari employee
	query := "SELECT id, identitynumber, name, employee_image_uri, gender, departement_id FROM employees WHERE 1=1"
	var args []interface{}
	argCount := 1

	// Menambahkan filter untuk identityNumber (prefix matching, case insensitive)
	if identityNumber != "" {
		query += fmt.Sprintf(" AND identitynumber ILIKE $%d", argCount)
		args = append(args, identityNumber+"%") // Pencarian berdasarkan prefix
		argCount++
	}

	// Menambahkan filter untuk name (prefix dan suffix matching, case insensitive)
	if name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argCount)
		args = append(args, "%"+name+"%") // Pencarian berdasarkan prefix atau suffix
		argCount++
	}

	// Menambahkan filter untuk gender
	if gender != "" {
		query += fmt.Sprintf(" AND gender = $%d", argCount)
		args = append(args, gender)
		argCount++
	}

	// Menambahkan filter untuk departmentId
	if departmentId != "0" { // 0 berarti departmentId tidak diberikan
		query += fmt.Sprintf(" AND department_id = $%d", argCount)
		args = append(args, departmentId)
		argCount++
	}

	// Menambahkan limit dan offset untuk pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	// Eksekusi query untuk mendapatkan data employees
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
