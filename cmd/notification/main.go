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
	"otus/internal/notification"
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
	notificationStore := notification.NewNotificationStore(db)
	notificationService := notification.NewNotificationService(rmq, *notificationStore)
	notificationService.ConsumeNotification("notification", "notification_service")
	notificationHandler := handlers.NewNotificationHandler(*notificationStore)
	metric.RegisterMetrics()
	r := mux.NewRouter()
	r.HandleFunc("/notification", metric.HttpMetricMiddleware(auth.AuthMiddleware(notificationHandler.GetNotifications), "/notification")).Methods("GET")
	r.HandleFunc("/notification/{OrderID}", metric.HttpMetricMiddleware(auth.AuthMiddleware(notificationHandler.GetNotificationByOrderID), "/notification/{OrderID}")).Methods("GET")
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")
	port := os.Getenv("SERVER_PORT")
	log.Println("Server running on " + port)
	log.Fatal(http.ListenAndServe(port, r))
}
