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

	existingUser, err := s.userRepo.GetByEmail(tx, ctx, email)
	if err != nil {
		log.WithError(err).Error("[service - Register]: Error checking for existing user")

		tx.Rollback()
		return nil, nil, err
	}

	if existingUser != nil {
		log.Infof("[service - Register]: User with email %s already exists", email)

		tx.Rollback()
		return nil, nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		log.WithError(err).Error("[service - Register]: Error hashing password")
		return nil, nil, err
	}

	user := &domain.User{
		FullName: fullname,
		Email:    email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(tx, ctx, user); err != nil {
		log.WithError(err).Error("[service - Register]: Failed to create user")

		tx.Rollback()
		return nil, nil, err
	}

	token, expiresAt, err := auth.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		log.WithError(err).Error("[service - Register]: Failed to generate token")

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
		log.WithError(err).Error("[service - Register]: Failed to create session")

		tx.Rollback()
		return nil, nil, err
	}

	if err := tx.Commit().Error; err != nil {
		log.WithError(err).Error("[service - Register]: Failed to commit transaction")
		return nil, nil, errors.New("internal server error")
	}

	return user, session, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*domain.User, *domain.Session, error) {
	log := logger.WithRequestID(ctx)

	user, err := s.userRepo.GetByEmail(s.db, ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Infof("[service - Login]: No user found with email %s", email)
			return nil, nil, errors.New("invalid email or password")
		}
		log.WithError(err).Error("[service - Login]: Error fetching user by email")
		return nil, nil, errors.New("internal server error")
	}

	if !auth.CheckPassword(password, user.Password) {
		log.Infof("[service - Login]: Invalid password for email %s", email)
		return nil, nil, errors.New("invalid email or password")
	}

	token, expiresAt, err := auth.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		log.WithError(err).Error("[service - Login]: Error generating token")
		return nil, nil, errors.New("internal server error")
	}

	session := &domain.Session{
		ID:           uuid.New(),
		SessionToken: token,
		UserID:       user.ID,
		ExpiresAt:    int(expiresAt.Unix()),
	}

	if err := s.sessionRepo.Create(s.db, ctx, session); err != nil {
		log.WithError(err).Error("[service - Login]: Error creating session")
		return nil, nil, err
	}

	return user, session, nil
}

func (s *authService) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	log := logger.WithRequestID(ctx)

	session, err := s.sessionRepo.GetByToken(s.db, ctx, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Infof("[service - GetUserByToken]: No session found with token %s", token)
			return nil, errors.New("invalid token")
		}
		log.WithError(err).Error("[service - GetUserByToken]: Error fetching session by token")
		return nil, errors.New("internal server error")
	}

	user, err := s.userRepo.GetByEmail(s.db, ctx, session.User.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Infof("[service - GetUserByToken]: No user found with ID %s", session.UserID)
			return nil, errors.New("user not found")
		}
		log.WithError(err).Error("[service - GetUserByToken]: Error fetching user by ID")
		return nil, errors.New("internal server error")
	}

	return user, nil
}
