package model

import (
	"database/sql"
	"nyala-backend/usecase/viewmodel"
	"time"
)

// orderDetailModel ...
type orderDetailModel struct {
	DB *sql.DB
}

// IOrderDetail ...
type IOrderDetail interface {
	Store(body viewmodel.OrderDetailVM, changedAt time.Time) (string, error)
}

// OrderDetailEntity ....
type OrderDetailEntity struct {
	OrderDetailID string         `db:"order_detail_id"`
	OrderID       string         `db:"order_id"`
	ProductID     string         `db:"product_id"`
	Qty           int            `db:"qty"`
	CreatedAt     string         `db:"created_at"`
	UpdatedAt     string         `db:"updated_at"`
	DeletedAt     sql.NullString `db:"deleted_at"`
}

// NewOrderDetailModel ...
func NewOrderDetailModel(db *sql.DB) IOrderDetail {
	return &orderDetailModel{DB: db}
}

// Store ...
func (model orderDetailModel) Store(body viewmodel.OrderDetailVM, changedAt time.Time) (res string, err error) {
	sql := `INSERT INTO "order_details" (
		"order_id", "product_id", "qty", "created_at", "updated_at"
		) VALUES($1, $2, $3, $4, $4) RETURNING "order_detail_id"`
	err = model.DB.QueryRow(sql, body.OrderID, body.ProductID, body.Qty, changedAt).Scan(&res)

	return res, err
}
