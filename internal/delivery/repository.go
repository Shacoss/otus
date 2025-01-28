package delivery

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

func NewDeliveryStore(db *sql.DB) *Store {
	return &Store{db: db, log: *logger.GetLogger()}
}

func (h *Store) CreateDelivery(delivery models.Delivery) (*int64, error) {
	var id int64
	err := h.db.QueryRow("INSERT INTO public.delivery (order_id, address, date) VALUES ($1, $2, $3) RETURNING id",
		delivery.OrderID, delivery.Address, delivery.Date.ToDateString()).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (h *Store) DeleteDeliveryByOrderID(orderID int64) error {
	_, err := h.db.Exec("DELETE FROM public.delivery WHERE order_id=$1", orderID)
	if err != nil {
		return err
	}
	return nil
}
