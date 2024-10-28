package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"otus/handlers"
)

var db *sql.DB

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}
	userHandler := &handlers.Handler{DB: db}
	router := mux.NewRouter()
	router.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	router.HandleFunc("/user", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/user/{id}", userHandler.GetUser).Methods("GET")
	router.HandleFunc("/user/{id}", userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/user/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Запуск сервера
	log.Println("Server starting on port 8000...")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
