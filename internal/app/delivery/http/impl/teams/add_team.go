package teams

import (
	"github.com/gofiber/fiber/v2"

	"pr-reviewer-assign-service/internal/app/delivery/http/api/teams"
	"pr-reviewer-assign-service/internal/app/domain/usecase"
	"pr-reviewer-assign-service/pkg/errors"
)

// AddTeam
//
//	@Summary	Создать команду с участниками (создаёт/обновляет пользователей)
//	@Tags		Teams
//	@Produce	json
//	@Param		body	body		teams.AddTeamParams	true	"teams.AddTeamParams"
//	@Success	200		{object}	teams.AddTeamResult
//	@Failure	400		{object}	api.ContractError
//	@Failure	404		{object}	api.ContractError
//	@Failure	500		{object}	api.ContractError
//	@Router		/team/add [post]
func (h *Handler) AddTeam(c *fiber.Ctx) error {
	var request teams.AddTeamParams

	err := c.BodyParser(&request)
	if err != nil {
		return errors.Wrap(err, errors.InternalError)
	}

	members := make([]usecase.TeamMemberParams, 0, len(request.Members))
	membersResult := make([]teams.AddTeamResultUser, 0, len(request.Members))
	for _, member := range request.Members {
		members = append(members, usecase.TeamMemberParams{
			UserID:   member.UserID,
			Username: member.UserName,
			IsActive: member.IsActive,
		})
		membersResult = append(membersResult, teams.AddTeamResultUser{
			UserID:   member.UserID,
			UserName: member.UserName,
			IsActive: member.IsActive,
		})
	}

	result, err := h.useCase.AddTeam(c.Context(), usecase.AddTeamParams{
		TeamName: request.TeamName,
		Members:  members,
	})
	if err != nil {
		return err
	}

	return c.JSON(teams.AddTeamResult{Team: teams.AddTeamResultTeam{
		TeamName: result.Team.TeamName,
		Members:  membersResult,
	}})
}
