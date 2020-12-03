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

// Store ...
func (model orderModel) Store(body viewmodel.OrderVM, changedAt time.Time) (res string, err error) {
	sql := `INSERT INTO "orders" (
		"customer_id", "order_number", "order_date", "payment_method_id", "created_at", "updated_at"
		) VALUES($1, $2, $3, $4, $5, $5) RETURNING "order_id"`
	err = model.DB.QueryRow(sql, body.CustomerID, body.OrderNumber, body.OrderDate, body.PaymentMethodID, changedAt).Scan(&res)

	return res, err
}
