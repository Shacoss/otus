package auth

import (
	"context"
	"net/http"
	"otus/pkg/logger"
	"strconv"
)

var log = logger.GetLogger()

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDString := r.Header.Get("X-UserID")
		if userIDString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userID, parseError := strconv.ParseInt(userIDString, 10, 64)
		if parseError != nil {
			log.Error(parseError.Error())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
