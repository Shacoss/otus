package billing

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

func NewBillingStore(db *sql.DB) *Store {
	return &Store{db: db, log: *logger.GetLogger()}
}

func (h *Store) CreateBilling(userID int64) (*int64, error) {
	var id int64
	err := h.db.QueryRow("INSERT INTO public.billing (user_id) VALUES ($1) RETURNING id", userID).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (h *Store) GetBillingByUserID(id int64) (*models.Billing, error) {
	var billing models.Billing
	err := h.db.QueryRow("SELECT * FROM public.billing WHERE user_id=$1", id).Scan(&billing.ID, &billing.UserID, &billing.Account)
	if err != nil {
		return nil, err
	}
	return &billing, nil
}

func (h *Store) UpdateBillingByUserID(billing models.Billing) error {
	_, err := h.db.Exec("UPDATE public.billing SET account=$1 WHERE user_id=$2", billing.Account, billing.UserID)
	if err != nil {
		return err
	}
	return nil
}
