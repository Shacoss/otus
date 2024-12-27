package handlers

import (
	"encoding/json"
	"net/http"
	"otus/internal/user"
	"otus/pkg/exception"
	models "otus/pkg/model"
)

type UserHandler struct {
	store user.Store
}

func NewUserHandler(store user.Store) *UserHandler {
	return &UserHandler{store: store}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	userProfile, err := h.store.GetUser(userID)
	isError := exception.HttpErrorHandler("User not found", err, w)
	if isError {
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(userProfile)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	userProfile, err := h.store.GetUser(userID)
	isError := exception.HttpErrorHandler("User not found", err, w)
	if isError {
		return
	}
	var userRequest models.User
	if parseError := json.NewDecoder(r.Body).Decode(&userRequest); parseError != nil {
		http.Error(w, "Failed to parse", http.StatusInternalServerError)
		return
	}
	if userRequest.Name == "" || userRequest.Email == "" {
		http.Error(w, "Incorrect user properties", http.StatusBadRequest)
		return
	}
	userUpdated, userUpdateError := h.store.UpdateUser(userProfile.ID, userRequest)
	if userUpdateError != nil {
		http.Error(w, "Failed to update user "+userRequest.Email, http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(userUpdated)
}
