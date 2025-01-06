package service

import (
	"context"
	"project_sprint/internal/utils"

	"project_sprint/internal/model"
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
