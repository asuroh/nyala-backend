package viewmodel

// UserVM ...
type UserVM struct {
	ID          string     `json:"id"`
	RoleID      string     `json:"role_id"`
	RoleName    string     `json:"role_name"`
	Data        string     `json:"data"`
	Information UserDataVM `json:"information"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	DeletedAt   string     `json:"deleted_at"`
}

// UserDataVM ...
type UserDataVM struct {
	Email    string   `json:"email"`
	Status   StatusVM `json:"status"`
	Password string   `json:"password"`
	UserName string   `json:"username"`
}

// StatusVM ...
type StatusVM struct {
	IsActive bool `json:"is_active"`
}

// AdminLoginVM ....
type AdminLoginVM struct {
	Token       string `json:"token"`
	ExpiredDate string `json:"expired_date"`
}
