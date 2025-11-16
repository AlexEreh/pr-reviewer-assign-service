package usecase

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"pr-reviewer-assign-service/internal/app/domain/model"
	"pr-reviewer-assign-service/pkg/log"
)

type SetIsActiveParams struct {
	UserID   string
	IsActive bool
}

type SetIsActiveResult struct {
	User model.User
}

func (u *UseCase) SetIsActive(
	ctx context.Context,
	params SetIsActiveParams,
) (SetIsActiveResult, error) {
	var result SetIsActiveResult

	err := u.txMan.Transactional(ctx, func(ctx context.Context) error {
		user, err := u.repo.GetUserByExternalID(ctx, params.UserID)
		if err != nil {
			log.LoggerFromCtx(ctx).
				Error("error getting user by external id",
					zap.Error(err),
					zap.String("UserID", params.UserID))

			return err
		}

		updatedUser, err := u.repo.SetUserActive(ctx, user.ID, params.IsActive)
		if err != nil {
			log.LoggerFromCtx(ctx).Error("failed to update user", zap.Error(err))

			return fmt.Errorf("failed to update user: %w", err)
		}

		teamMembers, err := u.repo.GetTeamMembersByUserID(ctx, user.ID)
		if err != nil {
			log.LoggerFromCtx(ctx).Error("error getting team members", zap.Error(err))

			return err
		}

		if len(teamMembers) == 0 {
			return fmt.Errorf("user has no team")
		}

		team, err := u.repo.GetTeamByID(ctx, teamMembers[0].TeamID)
		if err != nil {
			log.LoggerFromCtx(ctx).Error("error getting team by id", zap.Error(err))

			return fmt.Errorf("failed to get user team: %w", err)
		}

		result.User = model.User{
			UserID:   updatedUser.ExternalID,
			UserName: updatedUser.Username,
			TeamName: team.Name,
			IsActive: updatedUser.IsActive,
		}

		return nil
	})
	if err != nil {
		return SetIsActiveResult{}, err
	}

	return result, nil
}
