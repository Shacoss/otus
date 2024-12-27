package exception

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"net/http"
)

func HttpErrorHandler(message string, err error, w http.ResponseWriter) bool {
	if err == nil {
		return false
	}
	var pqErr *pq.Error
	switch {
	case errors.As(err, &pqErr):
		switch pqErr.Code {
		case "23505":
			http.Error(w, message, http.StatusConflict)
		default:
			http.Error(w, pqErr.Message, http.StatusInternalServerError)
		}
	case errors.Is(err, sql.ErrNoRows):
		http.Error(w, message, http.StatusNotFound)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return true
}
