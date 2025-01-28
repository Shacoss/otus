package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"otus/api/handlers"
	"otus/internal/auth"
	"otus/internal/warehouse"
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
	warehouseStore := warehouse.NewWarehouseStore(db)
	warehouseService := warehouse.NewWarehouseService(rmq, *warehouseStore)
	warehouseHandler := handlers.NewWarehouseHandler(*warehouseStore)
	warehouseService.ReservationProduct("reservation_product")
	warehouseService.ReservationRejectProduct("product_reservation_reject")
	r := mux.NewRouter()
	r.HandleFunc("/warehouse/product", auth.AuthMiddleware(warehouseHandler.GetAllProducts)).Methods("GET")
	r.HandleFunc("/warehouse/product/{ProductID}", auth.AuthMiddleware(warehouseHandler.GetProductByID)).Methods("GET")
	port := os.Getenv("SERVER_PORT")
	log.Println("Server running on " + port)
	log.Fatal(http.ListenAndServe(port, r))
}
