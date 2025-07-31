package main

import (
	"fmt"
	"net/http"
	"shekhar-patil/assignment/api/handlers"
	"shekhar-patil/assignment/api/middlewares"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	r.Handle("/pipeline", middlewares.AuthMiddleware(http.HandlerFunc(handlers.PipelineHandler))).Methods("POST")
	r.Handle("/report", middlewares.AuthMiddleware(http.HandlerFunc(handlers.ReportHandler))).Methods("GET")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	// Start server
	port := ":8080"
	fmt.Println("Server running on http://localhost" + port)
	http.ListenAndServe(port, r)
}
