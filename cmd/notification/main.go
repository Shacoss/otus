package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"otus/api/handlers"
	"otus/internal/auth"
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
	r := mux.NewRouter()
	r.HandleFunc("/notification", auth.AuthMiddleware(notificationHandler.GetNotifications)).Methods("GET")
	r.HandleFunc("/notification/{OrderID}", auth.AuthMiddleware(notificationHandler.GetNotificationByOrderID)).Methods("GET")
	port := os.Getenv("SERVER_PORT")
	log.Println("Server running on " + port)
	log.Fatal(http.ListenAndServe(port, r))
}
