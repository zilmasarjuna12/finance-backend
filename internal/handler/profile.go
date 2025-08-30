package handler

import (
	"finance-backend/internal/domain"
	"finance-backend/internal/model"
	"finance-backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type ProfileHandler struct {
	authService domain.AuthService
}

func NewProfileHandler(authService domain.AuthService) *ProfileHandler {
	return &ProfileHandler{
		authService: authService,
	}
}

func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	token := c.Locals("token").(string)

	user, err := h.authService.GetUserByToken(c.Context(), token)
	if err != nil {
		log.WithError(err).Error("[handler]: Failed to get user profile")
		return c.Status(fiber.StatusInternalServerError).JSON(
			model.NewResponseError("failed to get user profile"),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		model.NewResponseSuccess(
			model.User{
				ID:        user.ID.String(),
				Fullname:  user.FullName,
				Email:     user.Email,
				CreatedAt: int(user.CreatedAt),
				UpdatedAt: int(user.UpdatedAt),
			},
		),
	)
}
