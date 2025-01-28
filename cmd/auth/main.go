package main

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"otus/api/handlers"
	"otus/internal/auth"
	"otus/internal/user"
	"otus/pkg/broker"
	"otus/pkg/db"
)

var jwtSecret = []byte("supersecretkey")

func main() {
	db, err := db.CreateDbConnection()
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}
	defer db.Close()
	rmq := broker.NewRabbitMQ()
	err = rmq.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rmq.Close()
	userStore := user.NewUserStore(db)
	authService := auth.NewAuthService(rmq, "billing_account_create")
	authHandler := handlers.NewAuthHandler(*userStore, jwtSecret, *authService)
	r := mux.NewRouter()
	r.HandleFunc("/auth/register", authHandler.RegisterHandler).Methods("POST")
	r.HandleFunc("/auth/login", authHandler.LoginHandler).Methods("GET")
	r.HandleFunc("/auth/validate", authHandler.ValidateHandler).Methods("GET")
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	port := os.Getenv("SERVER_PORT")
	log.Println("Server running on " + port)
	log.Fatal(http.ListenAndServe(port, r))
}
