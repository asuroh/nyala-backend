package viewmodel

// CustomerVM ...
type CustomerVM struct {
	CustomerID   string `json:"customer_id"`
	CustomerName string `json:"customer_name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	Dob          string `json:"dob"`
	Sex          bool   `json:"sex"`
	Password     string `json:"password"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	DeletedAt    string `json:"deleted_at"`
}
