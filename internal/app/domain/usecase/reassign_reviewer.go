package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"pr-reviewer-assign-service/internal/app/data"
	"pr-reviewer-assign-service/internal/app/delivery/http/api"
	"pr-reviewer-assign-service/internal/app/domain/model"
	"pr-reviewer-assign-service/pkg/errors"
	"pr-reviewer-assign-service/pkg/log"
)

type ReassignReviewerParams struct {
	PullRequestID string
	OldReviewerID string
}

type ReassignReviewerResult struct {
	PR         model.PullRequest
	ReplacedBy string
}

func (u *UseCase) ReassignReviewer(
	ctx context.Context,
	params ReassignReviewerParams,
) (ReassignReviewerResult, error) {
	var result ReassignReviewerResult

	err := u.txMan.Transactional(ctx, func(ctx context.Context) error {
		pr, err := u.repo.GetPullRequestByExternalID(ctx, params.PullRequestID)
		if err != nil {
			return fmt.Errorf("PR not found")
		}

		if pr.Status == model.PullRequestStatusMerged {
			return fmt.Errorf("PR is merged")
		}

		oldReviewer, err := u.repo.GetUserByExternalID(ctx, params.OldReviewerID)
		if err != nil {
			return fmt.Errorf("reviewer not found")
		}

		reviewers, err := u.repo.GetCurrentReviewers(ctx, pr.ID)
		if err != nil {
			return fmt.Errorf("failed to get reviewers: %w", err)
		}

		var isAssigned bool
		var reviewerTeamID uuid.UUID
		for _, r := range reviewers {
			if r.ReviewerID == oldReviewer.ID {
				isAssigned = true
				reviewerTeamID = r.TeamID
				break
			}
		}

		if !isAssigned {
			return fmt.Errorf("reviewer not assigned")
		}

		newReviewer, err := u.findReplacementReviewer(ctx, reviewerTeamID, oldReviewer.ID)
		if err != nil {
			return fmt.Errorf("no candidate available")
		}

		err = u.replaceReviewer(ctx, pr.ID, oldReviewer.ID, newReviewer.ID, reviewerTeamID)
		if err != nil {
			return fmt.Errorf("failed to replace reviewer: %w", err)
		}

		history := data.PRReviewerHistory{
			ID:            uuid.New(),
			PullRequestID: pr.ID,
			OldReviewerID: sql.Null[data.UserInternalID]{V: oldReviewer.ID, Valid: true},
			NewReviewerID: newReviewer.ID,
			ChangedBy:     sql.Null[data.UserInternalID]{}, //nolint:exhaustruct // Системное изменение
			ChangedAt:     time.Now(),
			Reason:        model.PRReviewerHistoryChangeReasonReassignment,
		}

		_, err = u.repo.CreatePRReviewerHistory(ctx, history)
		if err != nil {
			log.LoggerFromCtx(ctx).Error("error creating review history", zap.Error(err))

			return fmt.Errorf("failed to log reassignment: %w", err)
		}

		updatedReviewers, err := u.repo.GetCurrentReviewers(ctx, pr.ID)
		if err != nil {
			log.LoggerFromCtx(ctx).Error("error getting current reviewers", zap.Error(err))

			return fmt.Errorf("failed to get updated reviewers: %w", err)
		}

		author, err := u.repo.GetUserByID(ctx, pr.AuthorID)
		if err != nil {
			log.LoggerFromCtx(ctx).Error("error getting author", zap.Error(err))

			return fmt.Errorf("failed to get author: %w", err)
		}

		var reviewerIDs []string

		for _, r := range updatedReviewers {
			user, err := u.repo.GetUserByID(ctx, r.ReviewerID)
			if err != nil {
				log.LoggerFromCtx(ctx).Error("error getting user by id", zap.Error(err))

				continue
			}

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
		result.ReplacedBy = newReviewer.ExternalID

		return nil
	})
	if err != nil {
		return ReassignReviewerResult{}, err
	}

	return result, nil
}

func (u *UseCase) findReplacementReviewer(
	ctx context.Context,
	teamID, excludeUserID uuid.UUID,
) (data.User, error) {
	teamMembers, err := u.repo.GetTeamMembersByTeamID(ctx, teamID)
	if err != nil {
		log.LoggerFromCtx(ctx).Error("error getting team members by team id",
			zap.Error(err),
			zap.String("TeamID", teamID.String()))

		return data.User{}, err
	}

	var availableUsers []data.User

	for _, tm := range teamMembers {
		user, err := u.repo.GetUserByID(ctx, tm.UserID)
		if err != nil {
			log.LoggerFromCtx(ctx).Error("error getting user by id",
				zap.Error(err),
				zap.String("UserID", tm.UserID.String()))

			continue
		}

		if user.IsActive && user.ID != excludeUserID {
			availableUsers = append(availableUsers, user)
		}
	}

	if len(availableUsers) == 0 {
		return data.User{}, errors.New(api.ErrNoCandidate)
	}

	selected := availableUsers[rand.Intn(len(availableUsers))] //nolint:gosec // Не принципиально

	return selected, nil
}

func (u *UseCase) replaceReviewer(
	ctx context.Context,
	prID, oldReviewerID, newReviewerID, teamID uuid.UUID,
) error {
	oldReviewers, err := u.repo.GetCurrentReviewers(ctx, prID)
	if err != nil {
		return err
	}

	for _, reviewer := range oldReviewers {
		if reviewer.ReviewerID == oldReviewerID {
			reviewer.IsCurrent = false

			now := time.Now()

			reviewer.ReplacedAt = sql.NullTime{
				Time:  now,
				Valid: true,
			}

			_, err := u.repo.UpdatePRReviewer(ctx, reviewer)
			if err != nil {
				log.LoggerFromCtx(ctx).Error("error updating pr reviewer",
					zap.Error(err),
					zap.Any("Reviewer", reviewer))

				return err
			}

			break
		}
	}

	newReviewer := data.PRReviewer{
		ID:            uuid.New(),
		PullRequestID: prID,
		ReviewerID:    newReviewerID,
		TeamID:        teamID,
		AssignedAt:    time.Now(),
		ReplacedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		IsCurrent:     true,
	}

	_, err = u.repo.CreatePRReviewer(ctx, newReviewer)
	if err != nil {
		return err
	}

	return nil
}
