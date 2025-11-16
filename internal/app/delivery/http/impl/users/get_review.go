package users

import (
	"github.com/gofiber/fiber/v2"

	"pr-reviewer-assign-service/internal/app/delivery/http/api"
	"pr-reviewer-assign-service/internal/app/delivery/http/api/users"
	"pr-reviewer-assign-service/internal/app/domain/usecase"
	"pr-reviewer-assign-service/pkg/errors"
)

// GetReview gets pending for review PRs of user
//
//	@Summary	Получить PR'ы, где пользователь назначен ревьювером
//	@Tags		Users
//	@Produce	json
//	@Param		user_id	query		string	true	"User ID"
//	@Success	200		{object}	GetReviewPRsResult
//	@Failure	400		{object}	api.ContractError
//	@Failure	404		{object}	api.ContractError
//	@Failure	500		{object}	api.ContractError
//	@Router		/users/getReview [get]
func (h *Handler) GetReview(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return errors.New(api.ErrUserIDNotProvided)
	}

	result, err := h.useCase.GetReview(c.Context(), usecase.GetReviewParams{UserID: userID})
	if err != nil {
		return err
	}

	pullRequestsResult := make([]users.GetReviewPRsResultPR, 0, len(result.PullRequests))
	for _, pr := range result.PullRequests {
		pullRequestsResult = append(pullRequestsResult, users.GetReviewPRsResultPR{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
		})
	}

	err = c.JSON(users.GetReviewPRsResult{
		UserID:       result.UserID,
		PullRequests: pullRequestsResult,
	})
	if err != nil {
		return err
	}

	return nil
}
