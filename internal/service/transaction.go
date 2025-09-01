package service

import (
	"context"
	"finance-backend/internal/constant"
	"finance-backend/internal/domain"
	"finance-backend/internal/model"
	"finance-backend/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type transactionService struct {
	db *gorm.DB

	transactionRepo domain.TransactionRepository
	walletRepo      domain.WalletRepository
}

func NewTransactionService(db *gorm.DB, transactionRepo domain.TransactionRepository, walletRepo domain.WalletRepository) domain.TransactionService {
	return &transactionService{
		db:              db,
		transactionRepo: transactionRepo,
		walletRepo:      walletRepo,
	}
}

func (s *transactionService) Create(ctx context.Context, userId string, request *model.CreateTransactionRequest) (*domain.Transaction, error) {
	log := logger.WithRequestID(ctx)

	log.Info("[service - transaction - Create]: Creating transaction")
	tx := s.db.Begin()

	if request.Type == constant.TransactionTypeIncome {
		if err := s.walletRepo.IncreaseBalance(tx, ctx, request.WalletID, request.Amount); err != nil {
			log.WithError(err).Error("[service - transaction - IncreaseBalance]: Failed to increase wallet balance")
			tx.Rollback()
			return nil, err
		}
	} else if request.Type == constant.TransactionTypeExpense {
		if err := s.walletRepo.DecreaseBalance(tx, ctx, request.WalletID, request.Amount); err != nil {
			log.WithError(err).Error("[service - transaction - DecreaseBalance]: Failed to decrease wallet balance")
			tx.Rollback()
			return nil, err
		}
	}

	transaction := &domain.Transaction{
		Amount:          request.Amount,
		Type:            request.Type,
		TransactionDate: request.TransactionDate,
		Note:            request.Note,
		WalletID:        uuid.MustParse(request.WalletID),
	}

	if request.BudgetID != nil && *request.BudgetID != "" {
		budgetID := uuid.MustParse(*request.BudgetID)
		transaction.BudgetID = &budgetID
	}

	if err := s.transactionRepo.Create(tx, ctx, userId, transaction); err != nil {
		log.WithError(err).Error("[service - transaction - Create]: Failed to create transaction")

		tx.Rollback()
		return nil, err
	}

	transaction, err := s.transactionRepo.GetDetail(tx, ctx, userId, transaction.ID.String())

	if err != nil {
		log.WithError(err).Error("[service - transaction - GetDetail]: Failed to get transaction detail after creation")

		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return transaction, nil
}

func (s *transactionService) GetList(ctx context.Context, userId string) ([]*domain.Transaction, error) {
	log := logger.WithRequestID(ctx)

	transactions, err := s.transactionRepo.GetList(s.db, ctx, userId)
	if err != nil {
		log.WithError(err).Error("[service - transaction - GetList]: Failed to get transaction list")
		return nil, err
	}

	return transactions, nil
}
