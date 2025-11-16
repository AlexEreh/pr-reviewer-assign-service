package usecase

import (
	"context"
	"fmt"

	"pr-reviewer-assign-service/internal/app/domain/model"
)

type GetTeamParams struct {
	TeamName string
}

type GetTeamResult struct {
	Team model.Team
}

func (u *UseCase) GetTeam(ctx context.Context, params GetTeamParams) (GetTeamResult, error) {
	team, err := u.repo.GetTeamByName(ctx, params.TeamName)
	if err != nil {
		return GetTeamResult{}, err
	}

	teamMembers, err := u.repo.GetTeamMembersByTeamID(ctx, team.ID)
	if err != nil {
		return GetTeamResult{}, fmt.Errorf("failed to get team members: %w", err)
	}

	var members []model.TeamMember
	for _, tm := range teamMembers {
		user, err := u.repo.GetUserByID(ctx, tm.UserID)
		if err != nil {
			return GetTeamResult{}, fmt.Errorf("failed to get user: %w", err)
		}
		members = append(members, model.TeamMember{
			UserID:   user.ExternalID,
			Username: user.Username,
			IsActive: user.IsActive,
		})
	}

	result := GetTeamResult{
		Team: model.Team{
			TeamName: team.Name,
			Members:  members,
		},
	}

	return result, nil
}
