package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"otus/internal/notification"
	"otus/pkg/exception"
	"otus/pkg/logger"
	"strconv"
)

type NotificationHandler struct {
	store notification.Store
	log   slog.Logger
}

func NewNotificationHandler(store notification.Store) *NotificationHandler {
	return &NotificationHandler{store: store, log: *logger.GetLogger()}
}

func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	notifications, err := h.store.GetNotificationByUserID(userID)
	isError := exception.HttpErrorHandler("Notifications not found", err, w)
	if isError {
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func (h *NotificationHandler) GetNotificationByOrderID(w http.ResponseWriter, r *http.Request) {
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
	userID := r.Context().Value("userID").(int64)
	orderNotification, err := h.store.GetNotificationByUserIDAndOrderID(userID, orderID)
	isError := exception.HttpErrorHandler("Notification not found by OrderID "+orderIDString, err, w)
	if isError {
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(orderNotification)
}
