package model

type User struct {
	ID        string `json:"id"`
	Fullname  string `json:"full_name"`
	Email     string `json:"email"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}
