package user

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

func NewUserStore(db *sql.DB) *Store {
	return &Store{db: db, log: *logger.GetLogger()}
}

func (h *Store) CreateUser(user models.User) (*int64, error) {
	var id int64
	err := h.db.QueryRow("INSERT INTO public.user (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		h.log.Error(err.Error())
		return nil, err
	}
	return &id, nil
}

func (h *Store) GetUser(id int64) (*models.User, error) {
	var user models.User
	err := h.db.QueryRow("SELECT id, name, email FROM public.user WHERE id=$1", id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		h.log.Error(err.Error())
		return nil, err
	}
	return &user, nil
}

func (h *Store) GetUserByEmailAndPassword(email string, password string) (*models.User, error) {
	var user models.User
	err := h.db.QueryRow("SELECT * FROM public.user WHERE email=$1 and password=$2", email, password).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		h.log.Error(err.Error())
		return nil, err
	}
	return &user, nil
}

func (h *Store) UpdateUser(id int64, user models.User) (*models.User, error) {
	_, err := h.db.Exec("UPDATE public.user SET name=$1, email=$2 WHERE id=$3", user.Name, user.Email, id)
	if err != nil {
		h.log.Error(err.Error())
		return nil, err
	}
	user.ID = id
	return &user, nil
}

func (h *Store) DeleteUser(id int64) error {
	_, err := h.db.Exec("DELETE FROM public.user WHERE id=$1", id)
	if err != nil {
		h.log.Error(err.Error())
		return err
	}
	return nil
}
