package statistics

import (
	"github.com/gofiber/fiber/v2"

	"pr-reviewer-assign-service/internal/app/delivery/http/api/statistics"
	"pr-reviewer-assign-service/internal/app/domain/usecase"
)

// GetStatistics
//
//	@Summary	Получить статистику по PR-ам
//	@Tags		Statistics
//	@Produce	json
//	@Success	200	{object}	statistics.GetStatisticsResult
//	@Failure	400	{object}	api.ContractError
//	@Failure	404	{object}	api.ContractError
//	@Failure	500	{object}	api.ContractError
//	@Router		/statistics/get [get]
func (h *Handler) GetStatistics(c *fiber.Ctx) error {
	result, err := h.useCase.GetStatistics(c.Context(), usecase.GetStatisticsParams{})
	if err != nil {
		return err
	}

	response := statistics.GetStatisticsResult{
		TotalPRs:        result.TotalPRs,
		OpenPRs:         result.OpenPRs,
		MergedPRs:       result.MergedPRs,
		UserAssignments: make([]statistics.UserAssignmentStats, 0, len(result.UserAssignments)),
		TeamStats:       make([]statistics.TeamStatistics, 0, len(result.TeamStats)),
		ReviewerLoad:    make([]statistics.ReviewerLoadStats, 0, len(result.ReviewerLoad)),
	}

	for _, stats := range result.UserAssignments {
		response.UserAssignments = append(response.UserAssignments, statistics.UserAssignmentStats{
			UserID:             stats.UserID,
			Username:           stats.Username,
			TeamName:           stats.TeamName,
			TotalPRs:           stats.TotalPRs,
			AssignedAsReviewer: stats.AssignedAsReviewer,
			ActiveAssignments:  stats.ActiveAssignments,
		})
	}

	for _, stats := range result.TeamStats {
		response.TeamStats = append(response.TeamStats, statistics.TeamStatistics{
			TeamName:     stats.TeamName,
			TotalPRs:     stats.TotalPRs,
			OpenPRs:      stats.OpenPRs,
			TotalReviews: stats.TotalReviews,
		})
	}

	for _, stats := range result.ReviewerLoad {
		response.ReviewerLoad = append(response.ReviewerLoad, statistics.ReviewerLoadStats{
			UserID:   stats.UserID,
			Username: stats.Username,
			Load:     stats.Load,
		})
	}

	err = c.JSON(response)
	if err != nil {
		return err
	}

	return nil
}
