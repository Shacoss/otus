package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"otus/api/handlers"
	"otus/internal/auth"
	"otus/internal/order"
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
	orderStore := order.NewOrderStore(db)
	orderService := order.NewOrderService(rmq, *orderStore)
	orderHandler := handlers.NewOrderHandler(*orderStore, *orderService)
	orderService.GetPaymentResult("oder_billing_response")
	r := mux.NewRouter()
	r.HandleFunc("/order", auth.AuthMiddleware(orderHandler.CreateOrder)).Methods("POST")
	r.HandleFunc("/order/{OrderID}", auth.AuthMiddleware(orderHandler.GetOrderByID)).Methods("GET")
	port := os.Getenv("SERVER_PORT")
	log.Println("Server running on " + port)
	log.Fatal(http.ListenAndServe(port, r))
}
