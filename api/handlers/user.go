package handlers

import (
	"encoding/json"
	"net/http"
	"otus/internal/user"
	models "otus/pkg/model"
	"strconv"
)

type UserHandler struct {
	store user.Store
}

func NewUserHandler(store user.Store) *UserHandler {
	return &UserHandler{store: store}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userProfile, err := getUser(h.store, w, r)
	if err == true {
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(userProfile)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userProfile, userProfileError := getUser(h.store, w, r)
	if userProfileError == true {
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

func getUser(db user.Store, w http.ResponseWriter, r *http.Request) (user *models.User, err bool) {
	userIDString := r.Header.Get("X-UserID")
	email := r.Header.Get("X-Email")
	if userIDString == "" || email == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil, true
	}
	userID, parseError := strconv.ParseInt(userIDString, 10, 64)
	if parseError != nil {
		http.Error(w, parseError.Error(), http.StatusInternalServerError)
		return nil, true
	}
	userProfile, userError := db.GetUser(userID)
	if userError != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return nil, true
	}
	if email != userProfile.Email {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil, true
	}
	return userProfile, false
}
