package order

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

func NewOrderStore(db *sql.DB) *Store {
	return &Store{db: db, log: *logger.GetLogger()}
}

func (h *Store) CreateOrder(order models.Order) (*int64, error) {
	var id int64
	err := h.db.QueryRow("INSERT INTO public.order (user_id, price, status) VALUES ($1, $2, $3) RETURNING id", order.UserID, order.Price, order.Status.String()).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (h *Store) GetOrderByID(id int64) (*models.Order, error) {
	var order models.Order
	var statusString string
	err := h.db.QueryRow("SELECT * FROM public.order WHERE id=$1", id).Scan(&order.ID, &order.UserID, &order.Price, &statusString)
	if err != nil {
		return nil, err
	}
	status, statusErr := models.ParseOrderStatus(statusString)
	if statusErr != nil {
		return nil, statusErr
	}
	order.Status = status
	return &order, nil
}

func (h *Store) UpdateOrderByID(order models.Order) error {
	_, err := h.db.Exec("UPDATE public.order SET status=$1 WHERE id=$2", order.Status.String(), order.ID)
	if err != nil {
		h.log.Error(err.Error())
		return err
	}
	return nil
}
