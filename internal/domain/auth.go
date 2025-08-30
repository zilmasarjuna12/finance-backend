package domain

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid()"`
	FullName string    `json:"full_name" gorm:"not null"`
	Email    string    `json:"email" gorm:"uniqueIndex;not null"`
	Password string    `json:"-" gorm:"not null"`

	CreatedAt int
	UpdatedAt int
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano"`

	Sessions []Session `gorm:"foreignKey:UserID"`
}

type Session struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid()"`
	SessionToken string    `gorm:"uniqueIndex;not null"`
	UserID       uuid.UUID `gorm:"type:uuid;not null"`
	ExpiresAt    int       `gorm:"not null"`

	User User `gorm:"foreignKey:UserID"`
}

type UserRepository interface {
	Create(db *gorm.DB, ctx context.Context, user *User) error
	GetByEmail(db *gorm.DB, ctx context.Context, email string) (*User, error)
}

type AuthService interface {
	Register(ctx context.Context, fullname, email, password string) (*User, *Session, error)
	Login(ctx context.Context, email, password string) (*User, *Session, error)
	GetUserByToken(ctx context.Context, token string) (*User, error)
}

type SessionRepository interface {
	Create(db *gorm.DB, ctx context.Context, session *Session) error
	GetByToken(db *gorm.DB, ctx context.Context, token string) (*Session, error)
	Delete(db *gorm.DB, ctx context.Context, token string) error
}

func (User) TableName() string {
	return "users"
}

func (Session) TableName() string {
	return "sessions"
}
