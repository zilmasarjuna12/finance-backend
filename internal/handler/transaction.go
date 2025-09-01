package handler

import (
	"finance-backend/internal/domain"
	"finance-backend/internal/model"
	"finance-backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	transactionService domain.TransactionService
}

func NewTransactionHandler(transactionService domain.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (h *TransactionHandler) Create(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	userId := c.Locals("userId").(string)

	var request model.CreateTransactionRequest
	if err := c.BodyParser(&request); err != nil {
		log.WithError(err).Error("[handler - transaction - Create]: Failed to parse create transaction request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	transaction, err := h.transactionService.Create(c.Context(), userId, &request)
	if err != nil {
		log.WithError(err).Error("[handler - transaction - Create]: Failed to create transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(model.NewResponseSuccess(model.Transaction{
		ID:              transaction.ID.String(),
		Amount:          transaction.Amount,
		Type:            transaction.Type,
		TransactionDate: transaction.TransactionDate,
		Note:            transaction.Note,
		Wallet: model.TransactionWallet{
			ID:   transaction.Wallet.ID.String(),
			Name: transaction.Wallet.Name,
		},
	}))
}

func (h *TransactionHandler) GetList(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	userId := c.Locals("userId").(string)

	transactions, err := h.transactionService.GetList(c.Context(), userId)
	if err != nil {
		log.WithError(err).Error("[handler - transaction - GetList]: Failed to get transaction list")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var response []model.Transaction
	for _, t := range transactions {
		tr := model.Transaction{
			ID:              t.ID.String(),
			Amount:          t.Amount,
			Type:            t.Type,
			TransactionDate: t.TransactionDate,
			Note:            t.Note,
			Wallet: model.TransactionWallet{
				ID:   t.Wallet.ID.String(),
				Name: t.Wallet.Name,
			},
		}

		if t.Budget != nil {
			tr.Budget = &model.TransactionBudget{
				ID:   t.Budget.ID.String(),
				Name: t.Budget.Name,
			}
		}

		response = append(response, tr)
	}

	return c.Status(fiber.StatusOK).JSON(model.NewResponseSuccess(response))
}
