package main

import (
	"family-tree-app/database"
	"family-tree-app/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database
	database.InitDB()

	r := mux.NewRouter()

	// Authentication routes
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("GET", "POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/upload", handlers.UploadHandler).Methods("GET", "POST")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("GET")
	r.HandleFunc("/tree/{id}", handlers.TreeHandler).Methods("GET")
	r.HandleFunc("/dashboard", handlers.DashboardHandler).Methods("GET")

	// Admin-only routes
	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(handlers.AuthMiddleware, handlers.AdminMiddleware) // Apply both middlewares

	admin.HandleFunc("/moderation", handlers.ModerationHandler).Methods("GET", "POST")

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":3002", r))
}
