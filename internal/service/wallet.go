package service

import (
	"context"
	"finance-backend/internal/domain"
	"finance-backend/internal/model"
	"finance-backend/pkg/logger"

	"gorm.io/gorm"
)

type walletService struct {
	db *gorm.DB

	walletRepo domain.WalletRepository
}

func NewWalletService(db *gorm.DB, walletRepo domain.WalletRepository) domain.WalletService {
	return &walletService{
		db:         db,
		walletRepo: walletRepo,
	}
}

func (s *walletService) Create(ctx context.Context, userId string, request *model.CreateWalletRequest) (*domain.Wallet, error) {
	log := logger.WithRequestID(ctx)

	tx := s.db.Begin()

	wallet := &domain.Wallet{
		Name:     request.Name,
		Type:     request.Type,
		Currency: request.Currency,
		Balance:  request.Balance,
	}

	if err := s.walletRepo.Create(tx, ctx, userId, wallet); err != nil {
		log.WithError(err).Error("[service - wallet - Create]: Failed to create wallet")
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return wallet, nil
}

func (s *walletService) GetList(ctx context.Context, userId string) ([]*domain.Wallet, error) {
	log := logger.WithRequestID(ctx)

	wallets, err := s.walletRepo.GetList(s.db, ctx, userId)
	if err != nil {
		log.WithError(err).Error("[service - wallet - GetList]: Failed to get wallet list")
		return nil, err
	}

	return wallets, nil
}
