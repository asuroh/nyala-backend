package request

// OrderRequest ....
type OrderRequest struct {
	PaymentMethodID string               `json:"payment_method_id" validate:"required"`
	OrderNumber     string               `json:"order_number"`
	OrderDetail     []OrderDetailRequest `json:"order_detail"`
}
