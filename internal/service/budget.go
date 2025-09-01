package service

import (
	"context"
	"finance-backend/internal/domain"
	"finance-backend/internal/model"
	"finance-backend/pkg/logger"

	"gorm.io/gorm"
)

type budgetService struct {
	db *gorm.DB

	budgetRepo domain.BudgetRepository
}

func NewBudgetService(db *gorm.DB, budgetRepo domain.BudgetRepository) domain.BudgetService {
	return &budgetService{
		db:         db,
		budgetRepo: budgetRepo,
	}
}

func (s *budgetService) Create(ctx context.Context, userId string, request *model.CreateBudgetRequest) (*domain.Budget, error) {
	log := logger.WithRequestID(ctx)

	tx := s.db.Begin()

	budget := &domain.Budget{
		Name:     request.Name,
		Amount:   request.Amount,
		Type:     request.Type,
		Category: request.Category,
	}

	if err := s.budgetRepo.Create(tx, ctx, userId, budget); err != nil {
		log.WithError(err).Error("[service - budget - Create]: Failed to create budget")
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return budget, nil
}

func (s *budgetService) GetList(ctx context.Context, userId string) ([]*domain.Budget, error) {
	log := logger.WithRequestID(ctx)

	budgets, err := s.budgetRepo.GetList(s.db, ctx, userId)
	if err != nil {
		log.WithError(err).Error("[service - budget - GetList]: Failed to get budget list")
		return nil, err
	}

	return budgets, nil
}
