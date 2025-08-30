package handler

import (
	"finance-backend/internal/domain"
	"finance-backend/internal/model"
	"finance-backend/pkg/logger"

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

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Error("[handler]: Failed to parse registration request body")
		return c.Status(fiber.StatusBadRequest).JSON(model.NewResponseError("invalid request body"))
	}

	user, session, err := h.authService.Register(c.Context(), req.Fullname, req.Email, req.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			return c.Status(fiber.StatusConflict).JSON(model.NewResponseError("user already exists"))
		}

		log.WithError(err).Error("[handler]: Failed to register user")
		return c.Status(fiber.StatusInternalServerError).JSON(model.NewResponseError("failed to register user"))
	}

	return c.Status(fiber.StatusCreated).JSON(model.NewResponseSuccess(
		model.AuthResponse{
			User: &model.User{
				ID:        user.ID.String(),
				Fullname:  user.FullName,
				Email:     user.Email,
				CreatedAt: int(user.CreatedAt),
				UpdatedAt: int(user.UpdatedAt),
			},
			Token:     session.SessionToken,
			ExpiresAt: session.ExpiresAt,
		},
	))
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Error("[handler]: Failed to parse login request body")
		return c.Status(fiber.StatusBadRequest).JSON(model.NewResponseError("invalid request body"))
	}

	user, session, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		if err.Error() == "invalid email or password" {
			return c.Status(fiber.StatusUnauthorized).JSON(model.NewResponseError("invalid email or password"))
		}

		log.WithError(err).Error("[handler]: Failed to login user")
		return c.Status(fiber.StatusInternalServerError).JSON(model.NewResponseError("failed to login user"))
	}

	return c.Status(fiber.StatusOK).JSON(model.NewResponseSuccess(
		model.AuthResponse{
			User: &model.User{
				ID:        user.ID.String(),
				Fullname:  user.FullName,
				Email:     user.Email,
				CreatedAt: int(user.CreatedAt),
				UpdatedAt: int(user.UpdatedAt),
			},
			Token:     session.SessionToken,
			ExpiresAt: session.ExpiresAt,
		}))
}
