package main

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"otus/api/handlers"
	"otus/internal/auth"
	"otus/internal/metric"
	"otus/internal/order"
	"otus/pkg/broker"
	"otus/pkg/db"
)

func main() {
	postgres, err := db.CreateDbConnection()
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}
	defer postgres.Close()
	rmq := broker.NewRabbitMQ()
	err = rmq.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rmq.Close()
	metric.RegisterMetrics()
	redis := db.CreateRedisClient()
	productURL := os.Getenv("PRODUCT_URL")
	orderStore := order.NewOrderStore(postgres)
	orderService := order.NewOrderService(rmq, *orderStore, *redis, productURL)
	orderHandler := handlers.NewOrderHandler(*orderStore, *orderService)
	orderService.ConsumeOrderResult("order_result")
	r := mux.NewRouter()
	r.HandleFunc("/order", metric.HttpMetricMiddleware(auth.AuthMiddleware(orderHandler.CreateOrder), "/order")).Methods("POST")
	r.HandleFunc("/order/{OrderID}", metric.HttpMetricMiddleware(auth.AuthMiddleware(orderHandler.GetOrderByID), "/order/{OrderID}")).Methods("GET")
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")
	port := os.Getenv("SERVER_PORT")
	log.Println("Server running on " + port)
	log.Fatal(http.ListenAndServe(port, r))
}
