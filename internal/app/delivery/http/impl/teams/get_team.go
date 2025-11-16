package teams

import (
	"github.com/gofiber/fiber/v2"

	"pr-reviewer-assign-service/internal/app/delivery/http/api"
	"pr-reviewer-assign-service/internal/app/delivery/http/api/teams"
	"pr-reviewer-assign-service/internal/app/domain/usecase"
	"pr-reviewer-assign-service/pkg/errors"
)

// GetTeam
//
//	@Summary	Получить команду с участниками
//	@Tags		Teams
//	@Produce	json
//	@Param		team_name	query		string	true	"Уникальное имя команды"
//	@Success	200			{object}	teams.GetTeamResult
//	@Failure	400			{object}	api.ContractError
//	@Failure	404			{object}	api.ContractError
//	@Failure	500			{object}	api.ContractError
//	@Router		/team/get [get]
func (h *Handler) GetTeam(c *fiber.Ctx) error {
	teamName := c.Query("team_name")
	if teamName == "" {
		return errors.New(api.ErrTeamNameNotProvided)
	}

	result, err := h.useCase.GetTeam(c.Context(), usecase.GetTeamParams{
		TeamName: teamName,
	})
	if err != nil {
		return err
	}

	membersResult := make([]teams.GetTeamResultUser, 0, len(result.Team.Members))
	for _, member := range result.Team.Members {
		membersResult = append(membersResult, teams.GetTeamResultUser{
			UserID:   member.UserID,
			UserName: member.Username,
			IsActive: member.IsActive,
		})
	}

	return c.JSON(teams.GetTeamResult{
		TeamName: result.Team.TeamName,
		Members:  membersResult,
	})
}
