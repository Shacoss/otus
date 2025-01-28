package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"otus/internal/warehouse"
	"otus/pkg/exception"
	"otus/pkg/logger"
	"strconv"
)

type WarehouseHandler struct {
	store warehouse.Store
	log   slog.Logger
}

func NewWarehouseHandler(store warehouse.Store) *WarehouseHandler {
	return &WarehouseHandler{store: store, log: *logger.GetLogger()}
}

func (h *WarehouseHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetAllProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *WarehouseHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productIDString := vars["ProductID"]
	if productIDString == "" {
		http.Error(w, "ProductID is empty", http.StatusBadRequest)
		return
	}
	productID, parseError := strconv.ParseInt(productIDString, 10, 64)
	if parseError != nil {
		http.Error(w, "Incorrect productID", http.StatusBadRequest)
		return
	}
	productByID, err := h.store.GetProductByID(productID)
	isError := exception.HttpErrorHandler("Product not found", err, w)
	if isError {
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(productByID)
}
