package repository

import (
	"context"
	"finance-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type transactionRepository struct{}

func NewTransactionRepository() domain.TransactionRepository {
	return &transactionRepository{}
}

func (r *transactionRepository) Create(db *gorm.DB, ctx context.Context, userId string, transaction *domain.Transaction) error {
	if err := db.WithContext(ctx).Create(transaction).Error; err != nil {
		return err
	}

	hasTransaction := domain.HasTransaction{
		UserID:        uuid.MustParse(userId),
		TransactionID: transaction.ID,
	}

	if err := db.WithContext(ctx).Create(&hasTransaction).Error; err != nil {
		return err
	}

	return nil
}

func (r *transactionRepository) GetList(db *gorm.DB, ctx context.Context, userId string) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction

	err := db.WithContext(ctx).
		Joins("JOIN has_transactions ON has_transactions.transaction_id = transactions.id").
		Where("has_transactions.user_id = ?", userId).
		Preload("Wallet").
		Preload("Budget").
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *transactionRepository) GetDetail(db *gorm.DB, ctx context.Context, userId string, transactionId string) (*domain.Transaction, error) {
	var transaction domain.Transaction

	err := db.WithContext(ctx).
		Joins("JOIN has_transactions ON has_transactions.transaction_id = transactions.id").
		Where("has_transactions.user_id = ? AND transactions.id = ?", userId, transactionId).
		Preload("Wallet").
		Preload("Budget").
		First(&transaction).Error
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
