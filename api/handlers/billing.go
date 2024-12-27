package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"otus/internal/billing"
	"otus/pkg/exception"
	"otus/pkg/logger"
	models "otus/pkg/model"
)

type BillingHandler struct {
	store billing.Store
	log   slog.Logger
}

func NewBillingHandler(store billing.Store) *BillingHandler {
	return &BillingHandler{store: store, log: *logger.GetLogger()}
}

func (h *BillingHandler) GetBilling(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	userBilling, err := h.store.GetBillingByUserID(userID)
	isError := exception.HttpErrorHandler("Billing not found", err, w)
	if isError {
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(userBilling)
}

func (h *BillingHandler) AddBillingAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	var billingRequest models.Billing
	if err := json.NewDecoder(r.Body).Decode(&billingRequest); err != nil {
		http.Error(w, "Failed to parse", http.StatusInternalServerError)
		return
	}
	userBilling, errGetBilling := h.store.GetBillingByUserID(userID)
	isError := exception.HttpErrorHandler("Billing not found", errGetBilling, w)
	if isError {
		return
	}
	userBilling.Account += billingRequest.Account
	errUpdateBilling := h.store.UpdateBillingByUserID(*userBilling)
	if errUpdateBilling != nil {
		http.Error(w, "Failed to add money", http.StatusInternalServerError)
		return
	}
}
