package usecase

import (
	"context"
	"fmt"

	"pr-reviewer-assign-service/internal/app/domain/model"
)

type GetReviewParams struct {
	UserID string
}

type GetReviewResult struct {
	UserID       string
	PullRequests []model.PullRequestShort
}

func (u *UseCase) GetReview(
	ctx context.Context,
	params GetReviewParams,
) (GetReviewResult, error) {
	user, err := u.repo.GetUserByExternalID(ctx, params.UserID)
	if err != nil {
		return GetReviewResult{}, fmt.Errorf("user not found")
	}

	reviewers, err := u.repo.GetUserAssignedPRs(ctx, user.ID)
	if err != nil {
		return GetReviewResult{}, fmt.Errorf("failed to get assigned PRs: %w", err)
	}

	var pullRequests []model.PullRequestShort

	for _, reviewer := range reviewers {
		pr, err := u.repo.GetPullRequestByID(ctx, reviewer.PullRequestID)
		if err != nil {
			continue
		}

		author, err := u.repo.GetUserByID(ctx, pr.AuthorID)
		if err != nil {
			continue
		}

		pullRequests = append(pullRequests, model.PullRequestShort{
			PullRequestID:   pr.ExternalID,
			PullRequestName: pr.Title,
			AuthorID:        author.ExternalID,
			Status:          pr.Status,
		})
	}

	result := GetReviewResult{
		UserID:       user.ExternalID,
		PullRequests: pullRequests,
	}

	return result, nil
}
