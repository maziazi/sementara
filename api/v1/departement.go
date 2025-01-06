package v1

import (
	"github.com/gorilla/mux"
	"net/http"
	"project_sprint/internal/handler"
)

func RegisterDepartmentRoutes(router *mux.Router) {
	router.HandleFunc("/departements", handler.CreateDepartmentHandler).Methods(http.MethodPost)
}
