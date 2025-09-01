package repository

import (
	"context"
	"finance-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type budgetRepository struct{}

func NewBudgetRepository() domain.BudgetRepository {
	return &budgetRepository{}
}

func (r *budgetRepository) Create(db *gorm.DB, ctx context.Context, userId string, budget *domain.Budget) error {
	if err := db.WithContext(ctx).Create(budget).Error; err != nil {
		return err
	}

	hasBudget := domain.HasBudget{
		UserID:   uuid.MustParse(userId),
		BudgetID: budget.ID,
	}

	if err := db.WithContext(ctx).Create(&hasBudget).Error; err != nil {
		return err
	}

	return nil
}

func (r *budgetRepository) GetList(db *gorm.DB, ctx context.Context, userId string) ([]*domain.Budget, error) {
	var budgets []*domain.Budget

	err := db.WithContext(ctx).
		Joins("JOIN has_budgets ON has_budgets.budget_id = budgets.id").
		Where("has_budgets.user_id = ?", userId).
		Find(&budgets).Error
	if err != nil {
		return nil, err
	}

	return budgets, nil
}
