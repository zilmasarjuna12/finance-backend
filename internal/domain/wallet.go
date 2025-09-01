package domain

import (
	"context"
	"finance-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type Wallet struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid()"`

	Name     string  `gorm:"type:varchar(100);not null"`
	Type     string  `gorm:"type:varchar(50);not null"`
	Currency string  `gorm:"type:varchar(10);not null"`
	Balance  float64 `gorm:"type:decimal(15,2);not null;check:balance >= 0"`

	CreatedAt int
	UpdatedAt int
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano"`
}

type HasWallet struct {
	UserID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	WalletID uuid.UUID `gorm:"type:uuid;primaryKey"`

	User   User   `gorm:"foreignKey:UserID;references:ID"`
	Wallet Wallet `gorm:"foreignKey:WalletID;references:ID"`
}

func (Wallet) TableName() string {
	return "wallets"
}

func (HasWallet) TableName() string {
	return "has_wallets"
}

type WalletRepository interface {
	Create(db *gorm.DB, ctx context.Context, userId string, wallet *Wallet) error
	GetList(db *gorm.DB, ctx context.Context, userId string) ([]*Wallet, error)
	DecreaseBalance(db *gorm.DB, ctx context.Context, walletId string, amount float64) error
	IncreaseBalance(db *gorm.DB, ctx context.Context, walletId string, amount float64) error
}

type WalletService interface {
	Create(ctx context.Context, userId string, request *model.CreateWalletRequest) (*Wallet, error)
	GetList(ctx context.Context, userId string) ([]*Wallet, error)
}
