package v1

import (
	"github.com/gorilla/mux"
	"net/http"
	"project_sprint/internal/handler"
)

func RegisterDepartmentRoutes(router *mux.Router) {
	routes := "/departments"
	router.HandleFunc(routes, handler.CreateDepartmentHandler).Methods(http.MethodPost)
	router.HandleFunc(routes, handler.GetDepartments).Methods(http.MethodGet)
}
