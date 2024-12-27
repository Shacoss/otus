package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"otus/api/handlers"
	"otus/internal/auth"
	"otus/internal/billing"
	"otus/pkg/broker"
	"otus/pkg/db"
)

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
	billingStore := billing.NewBillingStore(db)
	billingService := billing.NewBillingService(rmq, *billingStore)
	billingHandler := handlers.NewBillingHandler(*billingStore)
	billingService.CreateBilling("billing_create")
	billingService.RequestBilling("oder_billing_request", "oder_billing_response", "notification")
	r := mux.NewRouter()
	r.HandleFunc("/billing", auth.AuthMiddleware(billingHandler.GetBilling)).Methods("GET")
	r.HandleFunc("/billing/add", auth.AuthMiddleware(billingHandler.AddBillingAccount)).Methods("POST")
	port := os.Getenv("SERVER_PORT")
	log.Println("Server running on " + port)
	log.Fatal(http.ListenAndServe(port, r))
}
