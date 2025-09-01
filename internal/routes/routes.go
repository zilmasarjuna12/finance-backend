package routes

import (
	"finance-backend/internal/handler"
	"finance-backend/internal/repository"
	"finance-backend/internal/service"
	middleware "finance-backend/pkg/midleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	userRepository := repository.NewUserRepository()
	sessionRepository := repository.NewSessionRepository()
	walletRepository := repository.NewWalletRepository()
	budgetRepository := repository.NewBudgetRepository()
	transactionRepository := repository.NewTransactionRepository()

	authService := service.NewAuthService(db, userRepository, sessionRepository)
	walletService := service.NewWalletService(db, walletRepository)
	budgetService := service.NewBudgetService(db, budgetRepository)
	transactionService := service.NewTransactionService(db, transactionRepository, walletRepository)

	authHandler := handler.NewAuthHandler(authService)
	profileHandler := handler.NewProfileHandler(authService)
	walletHandler := handler.NewWalletHandler(walletService)
	budgetHandler := handler.NewBudgetHandler(budgetService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.LoggingMiddleware())
	app.Use(recover.New())

	// API v1 routes
	v1 := app.Group("/v1")

	// Health check route (public)
	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success":   true,
			"message":   "success",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Additional routes can be added here
	v1.Post("/auth/register", authHandler.Register)
	v1.Post("/auth/login", authHandler.Login)

	protected := v1.Group("/", middleware.AuthMiddleware(authService))

	protected.Get("/profile", profileHandler.GetProfile)

	protected.Get("/wallet", walletHandler.GetList)
	protected.Post("/wallet", walletHandler.Create)

	protected.Get("/budget", budgetHandler.GetList)
	protected.Post("/budget", budgetHandler.Create)

	protected.Post("/transaction", transactionHandler.Create)
	protected.Get("/transaction", transactionHandler.GetList)
}
