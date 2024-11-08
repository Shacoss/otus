package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"otus/handlers"
	"otus/metric"
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
	metric.RegisterMetrics()
	userHandler := &handlers.Handler{DB: db}
	router := mux.NewRouter()
	router.HandleFunc("/health", handlers.InstrumentedHandler(handlers.HealthHandler, "GET", "/health")).Methods("GET")
	router.HandleFunc("/user", handlers.InstrumentedHandler(userHandler.CreateUser, "POST", "/user")).Methods("POST")
	router.HandleFunc("/user/{id}", handlers.InstrumentedHandler(userHandler.GetUser, "GET", "/user/{id}")).Methods("GET")
	router.HandleFunc("/user/{id}", handlers.InstrumentedHandler(userHandler.UpdateUser, "PUT", "/user/{id}")).Methods("PUT")
	router.HandleFunc("/user/{id}", handlers.InstrumentedHandler(userHandler.DeleteUser, "DELETE", "/user/{id}")).Methods("DELETE")
	router.Handle("/metrics", promhttp.Handler())
	log.Println("Server starting on port 8000...")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
