package usecase

import (
	"context"
	"fmt"
	"time"

	"pr-reviewer-assign-service/internal/app/domain/model"
)

type MergePullRequestParams struct {
	PullRequestID string
}

type MergePullRequestResult struct {
	PR       model.PullRequest
	MergedAt time.Time
}

func (u *UseCase) MergePullRequest(
	ctx context.Context,
	params MergePullRequestParams,
) (MergePullRequestResult, error) {
	var result MergePullRequestResult

	err := u.txMan.Transactional(ctx, func(ctx context.Context) error {
		pr, err := u.repo.GetPullRequestByExternalID(ctx, params.PullRequestID)
		if err != nil {
			return fmt.Errorf("pr not found")
		}

		if pr.Status == model.PullRequestStatusMerged {
			author, _ := u.repo.GetUserByID(ctx, pr.AuthorID)
			reviewers, _ := u.repo.GetCurrentReviewers(ctx, pr.ID)

			var reviewerIDs []string

			for _, r := range reviewers {
				user, _ := u.repo.GetUserByID(ctx, r.ReviewerID)

				reviewerIDs = append(reviewerIDs, user.ExternalID)
			}

			result.PR = model.PullRequest{
				PullRequestShort: model.PullRequestShort{
					PullRequestID:   pr.ExternalID,
					PullRequestName: pr.Title,
					AuthorID:        author.ExternalID,
					Status:          pr.Status,
				},
				AssignedReviewers: reviewerIDs,
			}
			result.MergedAt = pr.MergedAt.Time

			return nil
		}

		mergedPR, err := u.repo.MergePullRequest(ctx, pr.ID)
		if err != nil {
			return fmt.Errorf("failed to merge PR: %w", err)
		}

		author, err := u.repo.GetUserByID(ctx, mergedPR.AuthorID)
		if err != nil {
			return fmt.Errorf("failed to get author: %w", err)
		}

		reviewers, err := u.repo.GetCurrentReviewers(ctx, mergedPR.ID)
		if err != nil {
			return fmt.Errorf("failed to get reviewers: %w", err)
		}

		var reviewerIDs []string

		for _, r := range reviewers {
			user, err := u.repo.GetUserByID(ctx, r.ReviewerID)
			if err != nil {
				continue
			}

			reviewerIDs = append(reviewerIDs, user.ExternalID)
		}

		result.PR = model.PullRequest{
			PullRequestShort: model.PullRequestShort{
				PullRequestID:   mergedPR.ExternalID,
				PullRequestName: mergedPR.Title,
				AuthorID:        author.ExternalID,
				Status:          mergedPR.Status,
			},
			AssignedReviewers: reviewerIDs,
		}
		result.MergedAt = mergedPR.MergedAt.Time

		return nil
	})
	if err != nil {
		return MergePullRequestResult{}, err
	}

	return result, nil
}
