package repository

import (
	"context"
	"finance-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type walletRepository struct{}

func NewWalletRepository() domain.WalletRepository {
	return &walletRepository{}
}

func (r *walletRepository) Create(db *gorm.DB, ctx context.Context, userId string, wallet *domain.Wallet) error {
	// create wallet
	if err := db.WithContext(ctx).Create(wallet).Error; err != nil {
		return err
	}

	hasWallet := domain.HasWallet{
		UserID:   uuid.MustParse(userId),
		WalletID: wallet.ID,
	}

	// create has_wallet
	if err := db.WithContext(ctx).Create(&hasWallet).Error; err != nil {
		return err
	}

	return nil
}

func (r *walletRepository) GetList(db *gorm.DB, ctx context.Context, userId string) ([]*domain.Wallet, error) {
	var wallets []*domain.Wallet

	err := db.WithContext(ctx).
		Joins("JOIN has_wallets ON has_wallets.wallet_id = wallets.id").
		Where("has_wallets.user_id = ?", userId).
		Find(&wallets).Error
	if err != nil {
		return nil, err
	}

	return wallets, nil
}

func (r *walletRepository) DecreaseBalance(db *gorm.DB, ctx context.Context, walletId string, amount float64) error {
	return db.WithContext(ctx).Model(&domain.Wallet{}).
		Where("id = ? AND balance >= ?", walletId, amount).
		Update("balance", gorm.Expr("balance - ?", amount)).Error
}

func (r *walletRepository) IncreaseBalance(db *gorm.DB, ctx context.Context, walletId string, amount float64) error {
	return db.WithContext(ctx).Model(&domain.Wallet{}).
		Where("id = ?", walletId).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}
