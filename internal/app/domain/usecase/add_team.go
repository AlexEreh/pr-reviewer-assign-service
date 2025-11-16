package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"pr-reviewer-assign-service/internal/app/data"
	"pr-reviewer-assign-service/internal/app/domain/model"
	"pr-reviewer-assign-service/pkg/log"
)

type AddTeamParams struct {
	TeamName string
	Members  []TeamMemberParams
}

type TeamMemberParams struct {
	UserID   string
	Username string
	IsActive bool
}

type AddTeamResult struct {
	Team model.Team
}

func (u *UseCase) AddTeam(ctx context.Context, params AddTeamParams) (AddTeamResult, error) {
	var result AddTeamResult

	err := u.txMan.Transactional(ctx, func(ctx context.Context) error {
		team := data.Team{
			ID:          uuid.New(),
			ExternalID:  uuid.NewString(),
			Name:        params.TeamName,
			Description: "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		createdTeam, err := u.repo.CreateTeam(ctx, team)
		if err != nil {
			log.LoggerFromCtx(ctx).Error("error creating team", zap.Error(err))

			return err
		}

		createdMembers := make([]model.TeamMember, 0, len(params.Members))

		for _, member := range params.Members {
			user := data.User{
				ID:         uuid.New(),
				ExternalID: member.UserID,
				Username:   member.Username,
				Email:      fmt.Sprintf("%s@example.com", member.UserID), // HACK Генерируем email
				IsActive:   member.IsActive,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			createdUser, err := u.repo.CreateUser(ctx, user)
			if err != nil {
				log.LoggerFromCtx(ctx).Error("error creating user", zap.Error(err))

				return err
			}

			teamMember := data.TeamMember{
				ID:        uuid.New(),
				TeamID:    createdTeam.ID,
				UserID:    createdUser.ID,
				Role:      "MEMBER",
				CreatedAt: time.Now(),
			}

			_, err = u.repo.CreateTeamMember(ctx, teamMember)
			if err != nil {
				log.LoggerFromCtx(ctx).Error("error creating team member", zap.Error(err))

				return err
			}

			createdMembers = append(createdMembers, model.TeamMember{
				UserID:   user.ExternalID,
				Username: user.Username,
				IsActive: user.IsActive,
			})
		}

		result.Team = model.Team{
			TeamName: createdTeam.Name,
			Members:  createdMembers,
		}

		return nil
	})
	if err != nil {
		return AddTeamResult{}, err
	}

	return result, nil
}
