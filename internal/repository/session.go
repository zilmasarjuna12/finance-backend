package repository

import (
	"context"
	"finance-backend/internal/domain"

	"gorm.io/gorm"
)

type sessionRepository struct {
}

func NewSessionRepository() domain.SessionRepository {
	return &sessionRepository{}
}

func (r *sessionRepository) Create(db *gorm.DB, ctx context.Context, session *domain.Session) error {
	return db.Create(session).Error
}

func (r *sessionRepository) GetByToken(db *gorm.DB, ctx context.Context, token string) (*domain.Session, error) {
	var session domain.Session
	err := db.WithContext(ctx).Preload("User").Where("session_token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) Delete(db *gorm.DB, ctx context.Context, token string) error {
	return db.WithContext(ctx).Where("session_token = ?", token).Delete(&domain.Session{}).Error
}
