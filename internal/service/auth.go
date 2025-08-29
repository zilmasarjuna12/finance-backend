package service

import (
	"context"
	"errors"
	"finance-backend/internal/domain"
	"finance-backend/pkg/auth"
	"finance-backend/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type authService struct {
	db *gorm.DB

	userRepo    domain.UserRepository
	sessionRepo domain.SessionRepository
}

func NewAuthService(db *gorm.DB, userRepo domain.UserRepository, sessionRepo domain.SessionRepository) domain.AuthService {
	return &authService{
		db:          db,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *authService) Register(ctx context.Context, fullname, email, password string) (*domain.User, *domain.Session, error) {
	log := logger.WithRequestID(ctx)

	tx := s.db.Begin()

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		log.WithError(err).Error("[service]: Error hashing password")
		return nil, nil, err
	}

	user := &domain.User{
		FullName: fullname,
		Email:    email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(tx, ctx, user); err != nil {
		log.WithError(err).Error("[service]: Failed to create user")

		tx.Rollback()
		return nil, nil, err
	}

	token, expiresAt, err := auth.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		log.WithError(err).Error("[service]: Failed to generate token")

		tx.Rollback()
		return nil, nil, err
	}

	session := &domain.Session{
		ID:           uuid.New(),
		SessionToken: token,
		UserID:       user.ID,
		ExpiresAt:    int(expiresAt.Unix()),
	}

	if err := s.sessionRepo.Create(tx, ctx, session); err != nil {
		log.WithError(err).Error("[service]: Failed to create session")

		tx.Rollback()
		return nil, nil, err
	}

	if err := tx.Commit().Error; err != nil {
		log.WithError(err).Error("[service]: Failed to commit transaction")
		return nil, nil, errors.New("internal server error")
	}

	return user, session, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*domain.User, *domain.Session, error) {
	log := logger.WithRequestID(ctx)

	user, err := s.userRepo.GetByEmail(s.db, ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Infof("[service]: No user found with email %s", email)
			return nil, nil, errors.New("invalid email or password")
		}
		log.WithError(err).Error("[service]: Error fetching user by email")
		return nil, nil, errors.New("internal server error")
	}

	if !auth.CheckPassword(password, user.Password) {
		log.Infof("[service]: Invalid password for email %s", email)
		return nil, nil, errors.New("invalid email or password")
	}

	token, expiresAt, err := auth.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		log.WithError(err).Error("[service]: Error generating token")
		return nil, nil, errors.New("internal server error")
	}

	session := &domain.Session{
		ID:           uuid.New(),
		SessionToken: token,
		UserID:       user.ID,
		ExpiresAt:    int(expiresAt.Unix()),
	}

	if err := s.sessionRepo.Create(s.db, ctx, session); err != nil {
		log.WithError(err).Error("[service]: Error creating session")
		return nil, nil, err
	}

	return user, session, nil
}
