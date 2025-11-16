package health

import (
	"github.com/gofiber/fiber/v2"

	"pr-reviewer-assign-service/internal/app/delivery/http/api/health"
	"pr-reviewer-assign-service/internal/app/domain/usecase"
)

// ReadyZ is an endpoint with checks on dependent services.
//
//	@Summary	Проверить, жив ли сам сервис + его зависимости (например, БД),
//
// Использовать, например, для readiness (startup?) probe куба.
//
//	@Tags		Health
//	@Produce	json
//	@Success	200
//	@Router		/health/readyz [get]
func (h *Handler) ReadyZ(c *fiber.Ctx) error {
	result, err := h.useCase.ReadyZ(c.Context(), usecase.ReadyZParams{})
	if err != nil {
		return err
	}

	return c.JSON(health.ReadyZResult{
		OK: result.OK,
	})
}
