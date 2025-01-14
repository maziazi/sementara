package service

import (
	"context"
	"fmt"
	"log"
	"project_sprint/internal/model"
	"project_sprint/internal/utils"
	"project_sprint/pkg/database"
	"strconv"
)

func CreateDepartment(dept model.Department) (model.Department, error) {

	if err := utils.ValidateDepartment(dept.Name); err != nil {
		return model.Department{}, err
	}
	query := `INSERT INTO department (name) VALUES ($1) RETURNING department_id`
	var deptID int

	err := database.GetDBPool().QueryRow(context.Background(), query, dept.Name).Scan(&deptID)
	if err != nil {
		return model.Department{}, err
	}

	return model.Department{
		DepartmentID: deptID,
		Name:         dept.Name,
	}, nil
}

func GetDepartments(limit, offset int, nameFilter string) ([]model.Department, error) {

	query := `SELECT department_id, name FROM department`
	var args []interface{}
	argIndex := 1

	if nameFilter != "" {
		query += ` WHERE name ILIKE $1`
		args = append(args, "%"+nameFilter+"%")
		argIndex++
	}

	query += ` LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := database.GetDBPool().Query(context.Background(), query, args...)
	if err != nil {
		log.Println("Error querying departments:", err)
		return nil, err
	}
	defer rows.Close()

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

func PatchDepartment(id int, dept model.Department) (*model.Department, error) {
	db := database.GetDBPool()

	query := `UPDATE department SET name = $1 WHERE department_id = $2 RETURNING department_id, name`
	row := db.QueryRow(context.Background(), query, dept.Name, id)

	var updatedDept model.Department
	err := row.Scan(&updatedDept.DepartmentID, &updatedDept.Name)
	if err != nil {
		return nil, err
	}

	return &updatedDept, nil
}

func DeleteDepartment(id int) error {
	// Cek apakah department ada di database
	var departmentName string
	err := database.GetDBPool().QueryRow(context.Background(), `SELECT name FROM department WHERE department_id = $1`, id).Scan(&departmentName)
	if err != nil {
		if err.Error() == "no rows in result set" {
			// Department tidak ditemukan
			return fmt.Errorf("department not found")
		}
		// Error lain saat query
		return err
	}

	// Cek apakah department masih memiliki employees
	var employeeCount int
	err = database.GetDBPool().QueryRow(context.Background(), `SELECT COUNT(*) FROM employees WHERE departement_id = $1`, id).Scan(&employeeCount)
	if err != nil {
		return err
	}

	// Jika masih ada employees yang terkait dengan department
	if employeeCount > 0 {
		return fmt.Errorf("department has employees")
	}

	// Hapus department dari database
	_, err = database.GetDBPool().Exec(context.Background(), `DELETE FROM department WHERE department_id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}
