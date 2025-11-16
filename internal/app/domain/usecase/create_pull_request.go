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

type CreatePullRequestParams struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
}

type CreatePullRequestResult struct {
	PR model.PullRequest
}

func (u *UseCase) CreatePullRequest(
	ctx context.Context,
	params CreatePullRequestParams,
) (CreatePullRequestResult, error) {
	var result CreatePullRequestResult

	err := u.txMan.Transactional(ctx, func(ctx context.Context) error {
		_, err := u.repo.GetPullRequestByExternalID(ctx, params.PullRequestID)
		if err == nil {
			return fmt.Errorf("PR already exists")
		}

		author, err := u.repo.GetUserByExternalID(ctx, params.AuthorID)
		if err != nil {
			log.LoggerFromCtx(ctx).Error("error getting user by external id", zap.Error(err))

			return errors.New(api.ErrNotFound)
		}

		teamMembers, err := u.repo.GetTeamMembersByUserID(ctx, author.ID)
		if err != nil {
			return fmt.Errorf("error getting team members by author id: %w", err)
		}

		if len(teamMembers) == 0 {
			return fmt.Errorf("author has no team")
		}

		team, err := u.repo.GetTeamByID(ctx, teamMembers[0].TeamID)
		if err != nil {
			log.LoggerFromCtx(ctx).Error("error getting team by ID", zap.Error(err))

			return fmt.Errorf("team not found")
		}

		pr := data.PullRequest{
			ID:                uuid.New(),
			ExternalID:        params.PullRequestID,
			Title:             params.PullRequestName,
			Description:       "",
			AuthorID:          author.ID,
			Status:            model.PullRequestStatusOpen,
			NeedMoreReviewers: false,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			MergedAt:          sql.NullTime{},
		}

		createdPR, err := u.repo.CreatePullRequest(ctx, pr)
		if err != nil {
			log.LoggerFromCtx(ctx).Error("error creating pr", zap.Error(err))

			return fmt.Errorf("failed to create PR: %w", err)
		}

		reviewers, err := u.assignReviewers(ctx, team.ID, author.ID)
		if err != nil {
			log.LoggerFromCtx(ctx).Error("error assigning reviewers", zap.Error(err))

			return fmt.Errorf("failed to assign reviewers: %w", err)
		}

		var assignedReviewerIDs []string

		for _, reviewer := range reviewers {
			prReviewer := data.PRReviewer{
				ID:            uuid.New(),
				PullRequestID: createdPR.ID,
				ReviewerID:    reviewer.ID,
				TeamID:        team.ID,
				AssignedAt:    time.Now(),
				ReplacedAt:    sql.NullTime{},
				IsCurrent:     true,
			}

			_, err := u.repo.CreatePRReviewer(ctx, prReviewer)
			if err != nil {
				log.LoggerFromCtx(ctx).
					Error("failed to assign reviewer",
						zap.Error(err),
						zap.Any("PRReviewer", prReviewer))

				return fmt.Errorf("failed to assign reviewer: %w", err)
			}

			assignedReviewerIDs = append(assignedReviewerIDs, reviewer.ExternalID)

			history := data.PRReviewerHistory{
				ID:            uuid.New(),
				PullRequestID: createdPR.ID,
				NewReviewerID: reviewer.ID,
				OldReviewerID: sql.Null[data.UserInternalID]{},
				ChangedBy:     sql.Null[data.UserInternalID]{V: author.ID, Valid: true},
				ChangedAt:     time.Now(),
				Reason:        model.PRReviewerHistoryChangeReasonInitial,
			}

			_, err = u.repo.CreatePRReviewerHistory(ctx, history)
			if err != nil {
				log.LoggerFromCtx(ctx).
					Error("failed to log reviewer assignment",
						zap.Error(err),
						zap.Any("History", history))

				return fmt.Errorf("failed to log reviewer assignment: %w", err)
			}
		}

		result.PR = model.PullRequest{
			PullRequestShort: model.PullRequestShort{
				PullRequestID:   createdPR.ExternalID,
				PullRequestName: createdPR.Title,
				AuthorID:        author.ExternalID,
				Status:          createdPR.Status,
			},
			AssignedReviewers: assignedReviewerIDs,
		}

		return nil
	})
	if err != nil {
		return CreatePullRequestResult{}, err
	}

	return result, nil
}

// assignReviewers выбирает до 2 случайных активных ревьюверов из команды (исключая автора)
func (u *UseCase) assignReviewers(
	ctx context.Context,
	teamID, authorID uuid.UUID,
) ([]data.User, error) {
	teamMembers, err := u.repo.GetTeamMembersByTeamID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	var availableReviewers []data.User

	for _, tm := range teamMembers {
		user, err := u.repo.GetUserByID(ctx, tm.UserID)
		if err != nil {
			continue
		}

		if user.IsActive && user.ID != authorID {
			availableReviewers = append(availableReviewers, user)
		}
	}

	count := min(2, len(availableReviewers))
	if count == 0 {
		return make([]data.User, 0), nil
	}

	rand.Shuffle(len(availableReviewers), func(i, j int) {
		availableReviewers[i], availableReviewers[j] = availableReviewers[j], availableReviewers[i]
	})

	return availableReviewers[:count], nil
}
