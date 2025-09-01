package handler

import (
	"finance-backend/internal/domain"
	"finance-backend/internal/model"
	"finance-backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type BudgetHandler struct {
	budgetService domain.BudgetService
}

func NewBudgetHandler(budgetService domain.BudgetService) *BudgetHandler {
	return &BudgetHandler{
		budgetService: budgetService,
	}
}

func (h *BudgetHandler) Create(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	userId := c.Locals("userId").(string)

	var request model.CreateBudgetRequest
	if err := c.BodyParser(&request); err != nil {
		log.WithError(err).Error("[handler - budget - Create]: Failed to parse create budget request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	budget, err := h.budgetService.Create(c.Context(), userId, &request)
	if err != nil {
		log.WithError(err).Error("[handler - budget - Create]: Failed to create budget")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(model.NewResponseSuccess(model.Budget{
		ID:        budget.ID.String(),
		Name:      budget.Name,
		Amount:    budget.Amount,
		Type:      budget.Type,
		Category:  budget.Category,
		CreatedAt: int(budget.CreatedAt),
		UpdatedAt: int(budget.UpdatedAt),
	}))
}

func (h *BudgetHandler) GetList(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	userId := c.Locals("userId").(string)

	budgets, err := h.budgetService.GetList(c.Context(), userId)
	if err != nil {
		log.WithError(err).Error("[handler - budget - GetList]: Failed to get budget list")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var response []model.Budget
	for _, budget := range budgets {
		response = append(response, model.Budget{
			ID:        budget.ID.String(),
			Name:      budget.Name,
			Amount:    budget.Amount,
			Type:      budget.Type,
			Category:  budget.Category,
			CreatedAt: int(budget.CreatedAt),
			UpdatedAt: int(budget.UpdatedAt),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.NewResponseSuccess(response))
}
