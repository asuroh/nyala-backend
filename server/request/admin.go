package request

// AdminRequest ....
type AdminRequest struct {
	Code           string `json:"code"`
	Name           string `json:"name" validate:"required"`
	Email          string `json:"email" validate:"email"`
	Password       string `json:"password"`
	RoleID         string `json:"role_id" validate:"required"`
	ProfileImageID string `json:"profile_image_id"`
	Status         bool   `json:"status"`
}

// UserRequest ...
type UserRequest struct {
	RoleID      string          `json:"role_id"`
	Information UserDataRequest `json:"information"`
}

// UserDataRequest ...
type UserDataRequest struct {
	Email    string        `json:"email"`
	Status   StatusRequest `json:"status"`
	Password string        `json:"password"`
	UserName string        `json:"username"`
}

// StatusRequest ...
type StatusRequest struct {
	IsActive bool `json:"is_active"`
}

// UserLoginRequest ....
type UserLoginRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
}
