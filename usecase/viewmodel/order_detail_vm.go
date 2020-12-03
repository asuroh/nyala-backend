package viewmodel

// OrderDetailVM ...
type OrderDetailVM struct {
	OrderDetailID string `json:"order_detail_id"`
	OrderID       string `json:"order_id"`
	ProductID     string `json:"product_id"`
	Qty           int    `json:"qty"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	DeletedAt     string `json:"deleted_at"`
}
