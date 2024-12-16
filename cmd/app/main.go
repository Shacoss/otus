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

func main() {
	db, err := util.CreateDbConnection()
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}
	defer db.Close()
	userStore := user.NewUserStore(db)
	userHandler := handlers.NewUserHandler(*userStore)
	r := mux.NewRouter()
	r.HandleFunc("/user/profile", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/user/profile", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	log.Println("Server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
