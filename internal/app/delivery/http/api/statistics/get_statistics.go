package statistics

import (
	"context"
)

type GetStatisticsParams struct{}

type GetStatisticsResult struct {
	TotalPRs        int64                 `json:"total_prs"`
	OpenPRs         int64                 `json:"open_prs"`
	MergedPRs       int64                 `json:"merged_prs"`
	UserAssignments []UserAssignmentStats `json:"user_assignments"`
	TeamStats       []TeamStatistics      `json:"team_stats"`
	ReviewerLoad    []ReviewerLoadStats   `json:"reviewer_load"`
}

type UserAssignmentStats struct {
	UserID             string  `json:"user_id"`
	Username           string  `json:"username"`
	TeamName           *string `json:"team_name"`
	TotalPRs           int64   `json:"total_prs"`
	AssignedAsReviewer int64   `json:"assigned_as_reviewer"`
	ActiveAssignments  int64   `json:"active_assignments"`
}

type TeamStatistics struct {
	TeamName     string `json:"team_name"`
	TotalPRs     int64  `json:"total_prs"`
	OpenPRs      int64  `json:"open_prs"`
	TotalReviews int64  `json:"total_reviews"`
}

type ReviewerLoadStats struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Load     int64  `json:"load"`
}

func (c Client) GetStatistics(
	ctx context.Context,
	params GetStatisticsParams,
) (GetStatisticsResult, error) {
	return GetStatisticsResult{}, nil
}
