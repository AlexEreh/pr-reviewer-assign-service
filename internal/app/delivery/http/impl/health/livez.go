package health

import "github.com/gofiber/fiber/v2"

// LiveZ is a "ping" endpoint with no checks.
//
//	@Summary	Проверить, жив ли исключительно сам сервис,
//
// Использовать, например, для liveness probe куба.
//
//	@Tags		Health
//	@Produce	json
//	@Success	200
//	@Router		/health/livez [get]
func (h *Handler) LiveZ(c *fiber.Ctx) error {
	c.Status(fiber.StatusOK)

	return nil
}
