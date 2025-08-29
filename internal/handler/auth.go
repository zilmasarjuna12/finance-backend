package handler

import (
	"finance-backend/internal/domain"
	"finance-backend/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Fullname string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	ID        string `json:"id"`
	Fullname  string `json:"full_name"`
	Email     string `json:"email"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

type AuthResponse struct {
	User      *User  `json:"user"`
	Token     string `json:"token"`
	ExpiresAt int    `json:"expires_at"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Error("[handler]: Failed to parse registration request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success":   false,
			"message":   "Invalid request body",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}

	user, session, err := h.authService.Register(c.Context(), req.Fullname, req.Email, req.Password)
	if err != nil {
		log.WithError(err).Error("[handler]: Failed to register user")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success":   false,
			"message":   "Failed to register user",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": AuthResponse{
			User: &User{
				ID:        user.ID.String(),
				Fullname:  user.FullName,
				Email:     user.Email,
				CreatedAt: int(user.CreatedAt),
				UpdatedAt: int(user.UpdatedAt),
			},
			Token:     session.SessionToken,
			ExpiresAt: session.ExpiresAt,
		},
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Error("[handler]: Failed to parse login request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success":   false,
			"message":   "Invalid request body",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}

	user, session, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		if err.Error() == "invalid email or password" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success":   false,
				"message":   "Invalid email or password",
				"timestamp": time.Now().Format(time.RFC3339),
			})
		}

		log.WithError(err).Error("[handler]: Failed to login user")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success":   false,
			"message":   err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login successful",
		"data": AuthResponse{
			User: &User{
				ID:        user.ID.String(),
				Fullname:  user.FullName,
				Email:     user.Email,
				CreatedAt: int(user.CreatedAt),
				UpdatedAt: int(user.UpdatedAt),
			},
			Token:     session.SessionToken,
			ExpiresAt: session.ExpiresAt,
		},
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
