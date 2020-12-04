package request

// OrderRequest ....
type OrderRequest struct {
	PaymentMethodID string               `json:"payment_method_id" validate:"required"`
	OrderDetail     []OrderDetailRequest `json:"order_detail"`
}
