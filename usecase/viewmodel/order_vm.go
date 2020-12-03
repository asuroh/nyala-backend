package viewmodel

// OrderVM ...
type OrderVM struct {
	OrderID         string          `json:"order_id"`
	CustomerID      string          `json:"customer_id"`
	OrderNumber     string          `json:"order_number"`
	OrderDate       string          `json:"order_date"`
	PaymentMethodID string          `json:"payment_method_id"`
	OrderDetail     []OrderDetailVM `json:"order_detail"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
	DeletedAt       string          `json:"deleted_at"`
}
