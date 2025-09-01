package domain

import (
	"context"
	"finance-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type Transaction struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid()"`

	Amount          float64 `gorm:"type:decimal(15,2);not null;check:amount > 0"`
	Type            string  `gorm:"type:varchar(50);not null"` // e.g., income, expense, transfer
	Note            string  `gorm:"type:varchar(255)"`
	TransactionDate int     `gorm:"not null"`

	CreatedAt int
	UpdatedAt int
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano"`

	WalletID uuid.UUID  `gorm:"type:uuid;not null"`
	BudgetID *uuid.UUID `gorm:"type:uuid"`

	Wallet Wallet  `gorm:"foreignKey:WalletID;references:ID"`
	Budget *Budget `gorm:"foreignKey:BudgetID;references:ID"`
}

type HasTransaction struct {
	UserID        uuid.UUID   `gorm:"type:uuid;primaryKey"`
	TransactionID uuid.UUID   `gorm:"type:uuid;primaryKey"`
	User          User        `gorm:"foreignKey:UserID;references:ID"`
	Transaction   Transaction `gorm:"foreignKey:TransactionID;references:ID"`
}

func (Transaction) TableName() string {
	return "transactions"
}

func (HasTransaction) TableName() string {
	return "has_transactions"
}

type TransactionRepository interface {
	Create(db *gorm.DB, ctx context.Context, userId string, transaction *Transaction) error
	GetDetail(db *gorm.DB, ctx context.Context, userId string, transactionId string) (*Transaction, error)
	GetList(db *gorm.DB, ctx context.Context, userId string) ([]*Transaction, error)
}

type TransactionService interface {
	Create(ctx context.Context, userId string, request *model.CreateTransactionRequest) (*Transaction, error)
	GetList(ctx context.Context, userId string) ([]*Transaction, error)
}
