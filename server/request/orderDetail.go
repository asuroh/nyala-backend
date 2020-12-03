package request

// OrderDetailRequest ....
type OrderDetailRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	OrderID   string `json:"order_detail"`
	Qty       int    `json:"qty" validate:"required"`
}
