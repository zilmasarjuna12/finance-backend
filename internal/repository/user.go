package repository

import (
	"context"
	"finance-backend/internal/domain"

	"gorm.io/gorm"
)

type userRepository struct {
}

func NewUserRepository() domain.UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(db *gorm.DB, ctx context.Context, user *domain.User) error {
	return db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByEmail(db *gorm.DB, ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
