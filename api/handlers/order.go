package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"otus/internal/order"
	"otus/pkg/exception"
	"otus/pkg/logger"
	models "otus/pkg/model"
	"strconv"
)

type OrderHandler struct {
	store   order.Store
	log     slog.Logger
	service order.Service
}

func NewOrderHandler(store order.Store, service order.Service) *OrderHandler {
	return &OrderHandler{store: store, log: *logger.GetLogger(), service: service}
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDString := vars["OrderID"]
	if orderIDString == "" {
		http.Error(w, "OrderID is empty", http.StatusBadRequest)
		return
	}
	orderID, parseError := strconv.ParseInt(orderIDString, 10, 64)
	if parseError != nil {
		http.Error(w, "Incorrect OrderID", http.StatusBadRequest)
		return
	}
	orderByID, err := h.store.GetOrderByID(orderID)
	isError := exception.HttpErrorHandler("Order not found", err, w)
	if isError {
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(orderByID)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	var orderRequest models.Order
	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		http.Error(w, "Failed to parse", http.StatusInternalServerError)
		return
	}
	validate := h.validateOrder(orderRequest)
	if validate != nil {
		http.Error(w, validate.Error(), http.StatusBadRequest)
		return
	}
	product, err := h.service.GetProductByID(orderRequest.Product.ID, userID)
	if err != nil {
		exception.HttpErrorHandler("Failed to create order. %s", err, w)
		return
	}
	orderRequest.Price = product.Price * float64(orderRequest.Product.Quantity)
	orderRequest.UserID = userID
	orderID, orderError := h.store.CreateOrder(orderRequest)
	if orderError != nil {
		exception.HttpErrorHandler("Failed to create order", orderError, w)
		return
	}
	orderRequest.ID = *orderID
	orderRequest.Delivery.OrderID = *orderID
	reservationProductErr := h.service.ReservationProduct(orderRequest)
	if reservationProductErr != nil {
		http.Error(w, "Failed to handle order", http.StatusInternalServerError)
		return
	}
	w.Header().Set("X-OrderID", strconv.FormatInt(*orderID, 10))
	w.WriteHeader(http.StatusAccepted)
}

func (h *OrderHandler) validateOrder(order models.Order) error {
	if order.Product.ID == 0 {
		return errors.New("product id is empty or 0")
	} else if order.Product.Quantity == 0 {
		return errors.New("product quantity is empty or 0")
	} else if order.Delivery.Date.IsZero() {
		return errors.New("delivery date is empty")
	} else if order.Delivery.Address == "" {
		return errors.New("delivery address is empty")
	}
	return nil
}
