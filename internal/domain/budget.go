package domain

import (
	"context"
	"finance-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type Budget struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid()"`

	Name     string  `gorm:"type:varchar(100);not null"`
	Amount   float64 `gorm:"type:decimal(15,2);not null;check:amount >= 0"`
	Type     string  `gorm:"type:varchar(50);not null"`
	Category string  `gorm:"type:varchar(50);not null"`

	CreatedAt int
	UpdatedAt int
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano"`
}

type HasBudget struct {
	UserID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	BudgetID uuid.UUID `gorm:"type:uuid;primaryKey"`

	User   User   `gorm:"foreignKey:UserID;references:ID"`
	Budget Budget `gorm:"foreignKey:BudgetID;references:ID"`
}

func (Budget) TableName() string {
	return "budgets"
}

func (HasBudget) TableName() string {
	return "has_budgets"
}

type BudgetRepository interface {
	Create(db *gorm.DB, ctx context.Context, userId string, budget *Budget) error
	GetList(db *gorm.DB, ctx context.Context, userId string) ([]*Budget, error)
}

type BudgetService interface {
	Create(ctx context.Context, userId string, request *model.CreateBudgetRequest) (*Budget, error)
	GetList(ctx context.Context, userId string) ([]*Budget, error)
}
