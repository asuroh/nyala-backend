package request

// CustomerRequest ....
type CustomerRequest struct {
	CustomerName string `json:"customer_name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	PhoneNumber  string `json:"phone_number" validate:"required"`
	Dob          string `json:"dob" validate:"required"`
	Sex          string `json:"sex" validate:"required"`
	SexBool      bool   `json:"sex_bool"`
	Password     string `json:"password" validate:"required"`
}

// CustomerLoginRequest ....
type CustomerLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
