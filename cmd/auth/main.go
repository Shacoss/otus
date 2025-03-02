package main

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"otus/api/handlers"
	"otus/internal/auth"
	"otus/internal/metric"
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
	metric.RegisterMetrics()
	r := mux.NewRouter()
	r.HandleFunc("/auth/register", metric.HttpMetricMiddleware(authHandler.RegisterHandler, "/auth/register")).Methods("POST")
	r.HandleFunc("/auth/login", metric.HttpMetricMiddleware(authHandler.LoginHandler, "/auth/login")).Methods("GET")
	r.HandleFunc("/auth/validate", metric.HttpMetricMiddleware(authHandler.ValidateHandler, "/auth/validate")).Methods("GET")
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")
	port := os.Getenv("SERVER_PORT")
	log.Println("Server running on " + port)
	log.Fatal(http.ListenAndServe(port, r))
}
