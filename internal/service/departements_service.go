package service

import (
	"context"
	"log"
	"project_sprint/internal/model"
	"project_sprint/internal/utils"
	"project_sprint/pkg/database"
)

func CreateDepartment(dept model.Department) (model.Department, error) {

	if err := utils.ValidateDepartment(dept.Name); err != nil {
		return model.Department{}, err
	}
	query := `INSERT INTO departments (name) VALUES ($1) RETURNING id`
	var deptID int

	err := database.DB.QueryRow(context.Background(), query, dept.Name).Scan(&deptID)
	if err != nil {
		return model.Department{}, err
	}

	return model.Department{
		DepartmentID: deptID,
		Name:         dept.Name,
	}, nil
}

func GetDepartments(limit, offset int, nameFilter string) ([]model.Department, error) {
	// Query dasar
	query := `SELECT id, name FROM departments`
	args := []interface{}{}

	// Jika ada filter nama, tambahkan kondisi WHERE
	if nameFilter != "" {
		query += ` WHERE LOWER(name) LIKE LOWER($1)`
		args = append(args, "%"+nameFilter+"%") // LIKE %abc% (prefix & suffix)
	}

	// Tambahkan limit & offset
	query += ` LIMIT $2 OFFSET $3`
	args = append(args, limit, offset)

	// Eksekusi query
	rows, err := database.DB.Query(context.Background(), query, args...)
	if err != nil {
		log.Println("Error querying departments:", err)
		return nil, err
	}
	defer rows.Close()

	// Parsing hasil query
	var departments []model.Department
	for rows.Next() {
		var dept model.Department
		if err := rows.Scan(&dept.DepartmentID, &dept.Name); err != nil {
			log.Println("Error scanning department:", err)
			return nil, err
		}
		departments = append(departments, dept)
	}

	return departments, nil
}

//func DeleteDepartment(dept model.Department) error {
//	query := `DELETE FROM departments WHERE id=$1`
//	var deptId int
//
//	err := database.DB.QueryRow(context.Background(), query, dept.DepartmentID).Scan(&deptId)
//
//	if err != nil {
//		return err
//	}
//	return model.Department{}
//}

//func GetDepartments(limit, offset int, nameFilter string) ([]model.Department, error) {
//	query := `SELECT id, name FROM departments WHERE name ILIKE $1 ORDER BY id LIMIT $2 OFFSET $3`
//	var departments []model.Department
//
//	// Tambahkan wildcard (%) untuk prefix dan suffix search
//	searchPattern := "%" + strings.ToLower(nameFilter) + "%"
//
//	rows, err := database.DB.Query(context.Background(), query, searchPattern, limit, offset)
//	if err != nil {
//		log.Println("Error querying departments:", err)
//		return nil, err
//	}
//	defer rows.Close()
//
//	// Loop hasil query
//	for rows.Next() {
//		var dept model.Department
//		if err := rows.Scan(&dept.DepartmentID, &dept.Name); err != nil {
//			log.Println("Error scanning department:", err)
//			return nil, err
//		}
//		departments = append(departments, dept)
//	}
//
//	return departments, nil
//}
