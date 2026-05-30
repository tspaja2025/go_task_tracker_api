package main

import (
	"context"
	"fmt"
	"log"
	"main/auth"
	"main/database"
	"main/middleware"
	"main/tasks"
	"net/http"
	"os"
)

func main() {
	log.Println("Starting API service...")

	// Initialize the database connection pool
	dbPool, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer dbPool.Close()

	// Initialize global ip rate limiter
	// 5 tokens per second refill, with maximum capacity of 10 tokens
	limiter := middleware.NewIPRateLimiter(5, 10)

	// Initialize Auth repo and handler
	authRepo := auth.NewRepository(dbPool)
	authHandler := auth.NewHandler(authRepo)

	// Initialize task repo and handler
	taskRepo := tasks.NewRepository(dbPool)
	taskHandler := tasks.NewHandler(taskRepo)

	// Routing
	http.HandleFunc("/register", limiter.Limit(authHandler.Register))
	http.HandleFunc("/login", limiter.Limit(authHandler.Login))
	http.HandleFunc("/refresh", limiter.Limit(authHandler.Refresh))

	// Protected route
	http.HandleFunc("/tasks/", limiter.Limit(middleware.AuthMiddleware(taskHandler.Router)))
	// http.HandleFunc("/tasks/test", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	// Get user ID
	// 	userID, _ := middleware.GetUserIDFromContext(r.Context())

	// 	w.Header().Set("Content-Type", "application/json")
	// 	fmt.Fprintf(w, `{"message": "Access granted!", "user_id:" %d}`, userID)
	// }))

	// Test route that quaries the database version
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		var version string
		// Query the database directly using the pool
		err := dbPool.QueryRow(context.Background(), "SELECT verion();").Scan(&version)
		if err != nil {
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `"status": "healthy", "database_version": "%s"`, version)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
