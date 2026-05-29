package main

import (
	"context"
	"fmt"
	"log"
	"main/database"
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
