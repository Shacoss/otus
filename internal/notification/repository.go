package notification

import (
	"database/sql"
	"log/slog"
	"otus/pkg/logger"
	models "otus/pkg/model"
)

type Store struct {
	db  *sql.DB
	log slog.Logger
}

func NewNotificationStore(db *sql.DB) *Store {
	return &Store{db: db, log: *logger.GetLogger()}
}

func (h *Store) CreateNotification(notification models.Notification) error {
	err := h.db.QueryRow("INSERT INTO public.notification (user_id, order_id, status) VALUES ($1, $2, $3)",
		notification.UserID, notification.OrderID, notification.Status.String())
	if err != nil {
		return err.Err()
	}
	return nil
}

func (h *Store) GetNotificationByUserID(id int64) (*models.Notification, error) {
	var notification models.Notification
	err := h.db.QueryRow("SELECT * FROM public.notification WHERE user_id=$1", id).Scan(&notification.UserID, &notification.OrderID, &notification.Status)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (h *Store) GetNotificationByUserIDAndOrderID(userID int64, orderID int64) (*models.Notification, error) {
	var notification models.Notification
	var statusString string
	err := h.db.QueryRow("SELECT * FROM public.notification WHERE user_id=$1 and order_id=$2", userID, orderID).Scan(&notification.UserID, &notification.OrderID, &statusString)
	if err != nil {
		return nil, err
	}
	status, statusErr := models.ParseOrderStatus(statusString)
	if statusErr != nil {
		return nil, statusErr
	}
	notification.Status = status
	return &notification, nil
}
