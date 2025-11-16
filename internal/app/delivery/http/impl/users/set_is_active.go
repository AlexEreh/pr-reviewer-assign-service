package users

import (
	"github.com/gofiber/fiber/v2"

	"pr-reviewer-assign-service/internal/app/delivery/http/api/users"
	"pr-reviewer-assign-service/internal/app/domain/usecase"
	"pr-reviewer-assign-service/pkg/errors"
)

// SetIsActive gets pending for review PRs of user
//
//	@Summary	Установить флаг активности пользователя
//	@Tags		Users
//	@Produce	json
//	@Param		body	body		users.SetIsActiveParams	true	"users.SetIsActiveParams"
//	@Success	200		{object}	users.SetIsActiveResult
//	@Failure	400		{object}	api.ContractError
//	@Failure	404		{object}	api.ContractError
//	@Failure	500		{object}	api.ContractError
//	@Router		/users/setIsActive [post]
func (h *Handler) SetIsActive(c *fiber.Ctx) error {
	var request users.SetIsActiveParams

	err := c.BodyParser(&request)
	if err != nil {
		return errors.Wrap(err, errors.InternalError)
	}

	result, err := h.useCase.SetIsActive(c.Context(), usecase.SetIsActiveParams{
		UserID:   request.UserID,
		IsActive: request.IsActive,
	})
	if err != nil {
		return err
	}

	err = c.JSON(users.SetIsActiveResult{
		User: users.SetIsActiveResultUser{
			UserID:   result.User.UserID,
			Username: result.User.UserName,
			TeamName: result.User.TeamName,
			IsActive: result.User.IsActive,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
