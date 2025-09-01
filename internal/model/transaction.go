package model

type CreateTransactionRequest struct {
	Amount          float64 `json:"amount"`
	Type            string  `json:"type"`
	Note            string  `json:"note"`
	TransactionDate int     `json:"transaction_date"`
	WalletID        string  `json:"wallet_id"`
	BudgetID        *string `json:"budget_id"`
}

type Transaction struct {
	ID              string             `json:"id"`
	Amount          float64            `json:"amount"`
	Type            string             `json:"type"`
	Note            string             `json:"note"`
	TransactionDate int                `json:"transaction_date"`
	Wallet          TransactionWallet  `json:"wallet"`
	Budget          *TransactionBudget `json:"budget"`
}

type TransactionWallet struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TransactionBudget struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
