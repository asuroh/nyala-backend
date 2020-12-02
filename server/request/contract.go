package request

// ForgotPasswordRequest ...
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// NewPasswordSubmitRequest ....
type NewPasswordSubmitRequest struct {
	Password string `json:"password" validate:"required,max=500"`
}
