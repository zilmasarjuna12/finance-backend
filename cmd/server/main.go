package main

import (
	"finance-backend/internal/routes"
	"finance-backend/pkg/database"
	"finance-backend/pkg/logger"
	"os"

	"finance-backend/pkg/migration"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	logger.InitLogger()
	log := logger.GetLogger()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Warn("No .env file found")
	}

	// Database connection
	dbConfig := database.GetConfigFromEnv()
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Info("Running database migrations...")
	sqlDB, err := database.GetSQLDB(db)
	if err != nil {
		log.Fatal("Failed to get SQL DB instance:", err)
	}

	migrationsDir := migration.GetMigrationsDir()
	if err := migration.MigrateUp(sqlDB, migrationsDir); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	log.Info("Database migrations completed successfully")

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log := logger.WithRequestID(c.Context())

			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			log.WithError(err).WithField("status_code", code).Error("Request error")

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Setup routes
	routes.SetupRoutes(app, db)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.WithField("port", port).Info("Starting server...")
	log.Fatal(app.Listen(":" + port))
}
