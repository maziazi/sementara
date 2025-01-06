package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project_sprint/internal/model"
	"project_sprint/internal/service"
	"project_sprint/internal/utils"
)

func CreateDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	var dept model.Department

	if err := json.NewDecoder(r.Body).Decode(&dept); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	departement, err := service.CreateDepartment(dept)
	err = utils.ValidateDepartment(departement.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	response := departement
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
