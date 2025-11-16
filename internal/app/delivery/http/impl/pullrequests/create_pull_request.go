package pullrequests

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"pr-reviewer-assign-service/internal/app/delivery/http/api/pullrequests"
	"pr-reviewer-assign-service/internal/app/domain/usecase"
)

// CreatePR
//
//	@Summary	Создать PR и автоматически назначить до 2 ревьюверов из команды автора
//	@Tags		PullRequests
//	@Produce	json
//	@Param		body	body		pullrequests.CreatePRParams	true	"pullrequests.CreatePRParams"
//	@Success	200		{object}	pullrequests.CreatePRResult
//	@Failure	400		{object}	api.ContractError
//	@Failure	404		{object}	api.ContractError
//	@Failure	500		{object}	api.ContractError
//	@Router		/pullRequest/create [post]
func (h *Handler) CreatePR(c *fiber.Ctx) error {
	var request pullrequests.CreatePRParams

	err := c.BodyParser(&request)
	if err != nil {
		return fmt.Errorf("error processing request body: %w", err)
	}

	result, err := h.useCase.CreatePullRequest(c.Context(), usecase.CreatePullRequestParams{
		PullRequestID:   request.PullRequestID,
		PullRequestName: request.PullRequestName,
		AuthorID:        request.AuthorID,
	})
	if err != nil {
		return err
	}

	return c.JSON(pullrequests.CreatePRResult{PR: pullrequests.CreatePRResultPR{
		PullRequestID:     result.PR.PullRequestID,
		PullRequestName:   result.PR.PullRequestName,
		AuthorID:          result.PR.AuthorID,
		Status:            result.PR.Status,
		AssignedReviewers: result.PR.AssignedReviewers,
	}})
}
