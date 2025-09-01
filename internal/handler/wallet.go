package handler

import (
	"finance-backend/internal/domain"
	"finance-backend/internal/model"
	"finance-backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type WalletHandler struct {
	walletService domain.WalletService
}

func NewWalletHandler(walletService domain.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

func (h *WalletHandler) Create(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	userId := c.Locals("userId").(string)

	var req model.CreateWalletRequest
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Error("[handler - wallet - Create]: Failed to parse create wallet request body")
		return c.Status(fiber.StatusBadRequest).JSON(model.NewResponseError("invalid request"))
	}

	wallet, err := h.walletService.Create(c.Context(), userId, &req)
	if err != nil {
		log.WithError(err).Error("[handler - wallet - Create]: Failed to create wallet")
		return c.Status(fiber.StatusInternalServerError).JSON(model.NewResponseError("failed to create wallet"))
	}

	return c.Status(fiber.StatusCreated).JSON(model.NewResponseSuccess(
		model.Wallet{
			ID:        wallet.ID.String(),
			Name:      wallet.Name,
			Type:      wallet.Type,
			Currency:  wallet.Currency,
			Balance:   wallet.Balance,
			CreatedAt: int(wallet.CreatedAt),
			UpdatedAt: int(wallet.UpdatedAt),
		},
	))
}

func (h *WalletHandler) GetList(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	userId := c.Locals("userId").(string)

	wallets, err := h.walletService.GetList(c.Context(), userId)
	if err != nil {
		log.WithError(err).Error("[handler - wallet - GetList]: Failed to get wallet list")
		return c.Status(fiber.StatusInternalServerError).JSON(model.NewResponseError("failed to get wallet list"))
	}

	var response []model.Wallet
	for _, wallet := range wallets {
		response = append(response, model.Wallet{
			ID:        wallet.ID.String(),
			Name:      wallet.Name,
			Type:      wallet.Type,
			Currency:  wallet.Currency,
			Balance:   wallet.Balance,
			CreatedAt: int(wallet.CreatedAt),
			UpdatedAt: int(wallet.UpdatedAt),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.NewResponseSuccess(response))
}
