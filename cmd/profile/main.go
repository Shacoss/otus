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
	"otus/pkg/db"
)

func main() {
	db, err := db.CreateDbConnection()
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}
	defer db.Close()
	userStore := user.NewUserStore(db)
	userHandler := handlers.NewUserHandler(*userStore)
	r := mux.NewRouter()
	metric.RegisterMetrics()
	r.HandleFunc("/user/profile", metric.HttpMetricMiddleware(auth.AuthMiddleware(userHandler.UpdateUser), "/user/profile")).Methods("PUT")
	r.HandleFunc("/user/profile", metric.HttpMetricMiddleware(auth.AuthMiddleware(userHandler.GetUser), "/user/profile")).Methods("GET")
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")
	port := os.Getenv("SERVER_PORT")
	log.Println("Server running on " + port)
	log.Fatal(http.ListenAndServe(port, r))
}
