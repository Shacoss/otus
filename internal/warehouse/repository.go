package warehouse

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

func NewWarehouseStore(db *sql.DB) *Store {
	return &Store{db: db, log: *logger.GetLogger()}
}

func (h *Store) GetAllProducts() ([]models.Product, error) {
	rows, err := h.db.Query("SELECT * FROM public.product")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products = make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (h *Store) GetProductByID(id int64) (*models.Product, error) {
	var product models.Product
	err := h.db.QueryRow("SELECT * FROM public.product WHERE id=$1", id).Scan(&product.ID, &product.Name, &product.Price, &product.Quantity)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (h *Store) UpdateProductQuantity(product models.Product) error {
	_, err := h.db.Exec("UPDATE public.product SET quantity=$1 WHERE id=$2", product.Quantity, product.ID)
	if err != nil {
		h.log.Error(err.Error())
		return err
	}
	return nil
}
