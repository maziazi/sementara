package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project_sprint/internal/model"
	"project_sprint/internal/service"
	"project_sprint/internal/utils"
	"strconv"
)

func CreateDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	var dept model.Department

	if err := json.NewDecoder(r.Body).Decode(&dept); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	department, err := service.CreateDepartment(dept)
	err = utils.ValidateDepartment(department.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	response := department
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetDepartments(w http.ResponseWriter, r *http.Request) {
	// Middleware auth (tambahkan jika ada)
	//token := r.Header.Get("Authorization")
	//if token == "" {
	//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//	return
	//}

	// Ambil query params
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	nameFilter := r.URL.Query().Get("name")

	// Konversi limit & offset ke integer, gunakan default jika tidak valid
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5 // Default limit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	// Panggil service untuk mendapatkan data department
	departments, err := service.GetDepartments(limit, offset, nameFilter)
	if err != nil {
		log.Println("Error fetching departments:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Kirim response JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(departments)
}
