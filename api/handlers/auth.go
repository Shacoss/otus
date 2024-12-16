package handlers

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"net/http"
	"otus/internal/user"
	models "otus/pkg/model"
	"strconv"
	"strings"
)

type AuthHandler struct {
	store     user.Store
	jwtSecret []byte
}

func NewAuthHandler(store user.Store, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{store: store, jwtSecret: jwtSecret}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	userRequest, isInvalid := getAuthUser(w, r)
	if isInvalid == true {
		return
	}
	userID, err := h.store.CreateUser(*userRequest)
	if err != nil && (err.(*pq.Error)).Code == "23505" {
		http.Error(w, "User already exist", http.StatusConflict)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userRequest.ID = *userID
	userRequest.Password = ""
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userRequest)
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	email, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	loginUser, err := h.store.GetUserByEmailAndPassword(email, password)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token, err := createJWT(*loginUser, h.jwtSecret)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) ValidateHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return h.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("X-UserID", claims.ID)
	w.Header().Set("X-Email", claims.Subject)
	w.WriteHeader(http.StatusOK)
}

func getAuthUser(w http.ResponseWriter, r *http.Request) (user *models.User, err bool) {
	var userRequest models.User
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		http.Error(w, "Failed to parse", http.StatusInternalServerError)
		return nil, true
	}
	if err := isValidUser(userRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, true
	}
	return &userRequest, false
}

func isValidUser(user models.User) error {
	if user.Name == "" {
		return errors.New("name is empty")
	}
	if user.Email == "" {
		return errors.New("email is empty")
	}
	if user.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}

func createJWT(user models.User, jwtSecret []byte) (string, error) {
	claims := &jwt.RegisteredClaims{
		ID:      strconv.FormatInt(user.ID, 10),
		Subject: user.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
