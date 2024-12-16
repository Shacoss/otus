package main

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"otus/api/handlers"
	"otus/internal/user"
	"otus/internal/util"
)

var jwtSecret = []byte("supersecretkey")

func main() {
	db, err := util.CreateDbConnection()
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}
	defer db.Close()
	userStore := user.NewUserStore(db)
	authHandler := handlers.NewAuthHandler(*userStore, jwtSecret)
	r := mux.NewRouter()
	r.HandleFunc("/auth/register", authHandler.RegisterHandler).Methods("POST")
	r.HandleFunc("/auth/login", authHandler.LoginHandler).Methods("GET")
	r.HandleFunc("/auth/validate", authHandler.ValidateHandler).Methods("GET")
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
