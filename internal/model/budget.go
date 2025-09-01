package model

type CreateBudgetRequest struct {
	Name     string  `json:"name" validate:"required"`
	Amount   float64 `json:"amount" validate:"required,numeric,min=0"`
	Type     string  `json:"type" validate:"required"`
	Category string  `json:"category" validate:"required"`
}

type Budget struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Amount    float64 `json:"amount"`
	Type      string  `json:"type"`
	Category  string  `json:"category"`
	CreatedAt int     `json:"created_at"`
	UpdatedAt int     `json:"updated_at"`
}
