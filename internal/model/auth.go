package model

type RegisterRequest struct {
	Fullname string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User      *User  `json:"user"`
	Token     string `json:"token"`
	ExpiresAt int    `json:"expires_at"`
}
