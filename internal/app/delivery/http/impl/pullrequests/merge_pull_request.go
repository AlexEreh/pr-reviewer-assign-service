package pullrequests

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"pr-reviewer-assign-service/internal/app/delivery/http/api/pullrequests"
	"pr-reviewer-assign-service/internal/app/domain/usecase"
)

// MergePR
//
//	@Summary	Пометить PR как MERGED (идемпотентная операция)
//	@Tags		PullRequests
//	@Produce	json
//	@Param		body	body		pullrequests.MergePRParams	true	"pullrequests.MergePRParams"
//	@Success	200		{object}	pullrequests.MergePRResult
//	@Failure	400		{object}	api.ContractError
//	@Failure	404		{object}	api.ContractError
//	@Failure	500		{object}	api.ContractError
//	@Router		/pullRequest/merge [post]
func (h *Handler) MergePR(c *fiber.Ctx) error {
	var request pullrequests.MergePRParams

	err := c.BodyParser(&request)
	if err != nil {
		return fmt.Errorf("error processing request body: %w", err)
	}

	result, err := h.useCase.MergePullRequest(c.Context(), usecase.MergePullRequestParams{
		PullRequestID: request.PullRequestID,
	})
	if err != nil {
		return err
	}

	return c.JSON(pullrequests.MergePRResult{PR: pullrequests.MergePRResultPR{
		PullRequestID:     result.PR.PullRequestID,
		PullRequestName:   result.PR.PullRequestName,
		AuthorID:          result.PR.AuthorID,
		Status:            result.PR.Status,
		AssignedReviewers: result.PR.AssignedReviewers,
		MergedAt:          result.MergedAt.Format(time.RFC3339),
	}})
}
