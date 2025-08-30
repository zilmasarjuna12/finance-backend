package middleware

import (
	"finance-backend/internal/domain"
	"finance-backend/pkg/auth"
	"finance-backend/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(authService domain.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := logger.WithRequestID(c.Context())

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Debug("[middleware - Auth]: Request attempted without authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success":   false,
				"message":   "Authorization header required",
				"timestamp": time.Now().Format(time.RFC3339),
			})
		}

		token, err := auth.ExtractTokenFromBearer(authHeader)
		if err != nil {
			log.WithError(err).Debug("[middleware - Auth]: Failed to extract token from authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success":   false,
				"error":     "Invalid authorization",
				"timestamp": time.Now().Format(time.RFC3339),
			})
		}

		user, err := authService.GetUserByToken(c.Context(), token)
		if err != nil {
			log.WithError(err).Debug("[middleware - Auth]: Failed to get user by token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success":   false,
				"error":     "Invalid token",
				"timestamp": time.Now().Format(time.RFC3339),
			})
		}

		c.Locals("user", user)
		c.Locals("token", token)

		log.WithField("user_id", user.ID).Debug("[middleware - Auth]: User authenticated successfully")

		return c.Next()
	}
}
