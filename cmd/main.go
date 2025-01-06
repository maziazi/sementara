package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"project_sprint/api/v1"
	"project_sprint/pkg/database"
)

func main() {

	database.InitDB()
	defer database.CloseDB()

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/v1").Subrouter()

	v1.RegisterDepartmentRoutes(apiRouter)

	log.Println("Server started on http://localhost:8080")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", router))
}
