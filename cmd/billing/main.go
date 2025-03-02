package main

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"otus/api/handlers"
	"otus/internal/auth"
	"otus/internal/billing"
	"otus/internal/metric"
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
	billingService.CreateBillingAccount("billing_account_create")
	billingService.CreatePayment("create_payment")
	metric.RegisterMetrics()
	r := mux.NewRouter()
	r.HandleFunc("/billing", metric.HttpMetricMiddleware(auth.AuthMiddleware(billingHandler.GetBilling), "/billing")).Methods("GET")
	r.HandleFunc("/billing/add", metric.HttpMetricMiddleware(auth.AuthMiddleware(billingHandler.AddBillingAccount), "/billing/add")).Methods("POST")
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")
	port := os.Getenv("SERVER_PORT")
	log.Println("Server running on " + port)
	log.Fatal(http.ListenAndServe(port, r))
}
