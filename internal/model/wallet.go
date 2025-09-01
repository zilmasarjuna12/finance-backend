package model

type CreateWalletRequest struct {
	Name     string  `json:"name" validate:"required"`
	Type     string  `json:"type" validate:"required,oneof=personal business"`
	Currency string  `json:"currency" validate:"required,len=3"`
	Balance  float64 `json:"balance" validate:"required,numeric"`
}

type Wallet struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Currency  string  `json:"currency"`
	Balance   float64 `json:"balance"`
	CreatedAt int     `json:"created_at"`
	UpdatedAt int     `json:"updated_at"`
}
