package model

import (
	"database/sql"
	"nyala-backend/usecase/viewmodel"
	"time"
)

// orderModel ...
type orderModel struct {
	DB *sql.DB
}

// IOrder ...
type IOrder interface {
	Count() (int, error)
	Store(body viewmodel.OrderVM, changedAt time.Time) (string, error)
}

// OrderEntity ....
type OrderEntity struct {
	OrderID         string         `db:"order_id"`
	CustomerID      string         `db:"customer_id"`
	OrderNumber     string         `db:"order_number"`
	OrderDate       string         `db:"order_date"`
	PaymentMethodID string         `db:"payment_method_id"`
	CreatedAt       string         `db:"created_at"`
	UpdatedAt       string         `db:"updated_at"`
	DeletedAt       sql.NullString `db:"deleted_at"`
}

// NewOrderModel ...
func NewOrderModel(db *sql.DB) IOrder {
	return &orderModel{DB: db}
}

// Count ...
func (model orderModel) Count() (count int, err error) {
	query := `SELECT COUNT(order_id) FROM orders WHERE deleted_at IS NULL AND EXTRACT(MONTH FROM order_date) = EXTRACT(MONTH FROM NOW()) AND EXTRACT(YEAR FROM order_date) = EXTRACT(YEAR FROM NOW())`
	err = model.DB.QueryRow(query).Scan(&count)

	return count, err
}

// Store ...
func (model orderModel) Store(body viewmodel.OrderVM, changedAt time.Time) (res string, err error) {
	sql := `INSERT INTO "orders" (
		"customer_id", "order_number", "order_date", "payment_method_id", "created_at", "updated_at"
		) VALUES($1, $2, $3, $4, $5, $5) RETURNING "order_id"`
	err = model.DB.QueryRow(sql, body.CustomerID, body.OrderNumber, body.OrderDate, body.PaymentMethodID, changedAt).Scan(&res)

	return res, err
}
