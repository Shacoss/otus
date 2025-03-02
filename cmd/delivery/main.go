package main

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"otus/internal/delivery"
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
	metric.RegisterMetrics()
	deliverStore := delivery.NewDeliveryStore(db)
	deliveryService := delivery.NewDeliveryService(rmq, *deliverStore)
	deliveryService.RejectDelivery("reject_delivery")
	deliveryService.CreateDelivery("create_delivery")
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")
	port := os.Getenv("SERVER_PORT")
	log.Println("Server running on " + port)
	log.Fatal(http.ListenAndServe(port, r))
	//quit := make(chan os.Signal, 1)
	//signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	//deliveryService.RejectDelivery("reject_delivery")
	//deliveryService.CreateDelivery("create_delivery")
	//<-quit
	//log.Println("Shutting down...")
}
