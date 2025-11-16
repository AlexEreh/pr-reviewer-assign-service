package pullrequests

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"pr-reviewer-assign-service/internal/app/delivery/http/api/pullrequests"
	"pr-reviewer-assign-service/internal/app/domain/usecase"
)

// ReassignPR
//
//	@Summary	Переназначить конкретного ревьювера на другого из его команды
//	@Tags		PullRequests
//	@Produce	json
//	@Param		body	body		pullrequests.ReassignPRParams	true	"pullrequests.ReassignPRParams"
//	@Success	200		{object}	pullrequests.ReassignPRResult
//	@Failure	400		{object}	api.ContractError
//	@Failure	404		{object}	api.ContractError
//	@Failure	409		{object}	api.ContractError
//	@Failure	500		{object}	api.ContractError
//	@Router		/pullRequest/reassign [post]
func (h *Handler) ReassignPR(c *fiber.Ctx) error {
	var request pullrequests.ReassignPRParams

	err := c.BodyParser(&request)
	if err != nil {
		return fmt.Errorf("error processing request body: %w", err)
	}

	result, err := h.useCase.ReassignReviewer(c.Context(), usecase.ReassignReviewerParams{
		PullRequestID: request.PullRequestID,
		OldReviewerID: request.OldReviewerID,
	})
	if err != nil {
		return err
	}

	return c.JSON(pullrequests.ReassignPRResult{
		PR: pullrequests.ReassignPRResultPR{
			PullRequestID:     result.PR.PullRequestID,
			PullRequestName:   result.PR.PullRequestName,
			AuthorID:          result.PR.AuthorID,
			Status:            result.PR.Status,
			AssignedReviewers: result.PR.AssignedReviewers,
		},
		ReplacedBy: result.ReplacedBy,
	})
}
