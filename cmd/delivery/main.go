package main

import (
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"otus/internal/delivery"
	"otus/pkg/broker"
	"otus/pkg/db"
	"syscall"
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
	deliverStore := delivery.NewDeliveryStore(db)
	deliveryService := delivery.NewDeliveryService(rmq, *deliverStore)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	deliveryService.RejectDelivery("reject_delivery")
	deliveryService.CreateDelivery("create_delivery")
	<-quit
	log.Println("Shutting down...")
}
