package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"pr-reviewer-assign-service/internal/app/data"
	"pr-reviewer-assign-service/internal/app/domain/model"
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

func (u *UseCase) GetStatistics(
	ctx context.Context,
	_ GetStatisticsParams,
) (GetStatisticsResult, error) {
	var result GetStatisticsResult

	allPRs, err := u.getAllPullRequests(ctx)
	if err != nil {
		return GetStatisticsResult{}, fmt.Errorf("failed to get PRs: %w", err)
	}

	result.TotalPRs = int64(len(allPRs))
	result.OpenPRs = u.countPRsByStatus(allPRs, model.PullRequestStatusOpen)
	result.MergedPRs = u.countPRsByStatus(allPRs, model.PullRequestStatusMerged)

	userStats, err := u.calculateUserStatistics(ctx)
	if err != nil {
		return GetStatisticsResult{}, fmt.Errorf("failed to calculate user stats: %w", err)
	}

	result.UserAssignments = userStats

	teamStats, err := u.calculateTeamStatistics(ctx, allPRs)
	if err != nil {
		return GetStatisticsResult{}, fmt.Errorf("failed to calculate team stats: %w", err)
	}

	result.TeamStats = teamStats

	reviewerLoad, err := u.calculateReviewerLoad(ctx)
	if err != nil {
		return GetStatisticsResult{}, fmt.Errorf("failed to calculate reviewer load: %w", err)
	}

	result.ReviewerLoad = reviewerLoad

	return result, nil
}

func (u *UseCase) getAllPullRequests(ctx context.Context) ([]data.PullRequest, error) {
	openPRs, err := u.repo.GetPullRequestsByStatus(ctx, model.PullRequestStatusOpen)
	if err != nil {
		return nil, err
	}

	mergedPRs, err := u.repo.GetPullRequestsByStatus(ctx, model.PullRequestStatusMerged)
	if err != nil {
		return nil, err
	}

	return append(openPRs, mergedPRs...), nil
}

func (u *UseCase) countPRsByStatus(prs []data.PullRequest, status string) int64 {
	var count int64

	for _, pr := range prs {
		if pr.Status == status {
			count++
		}
	}

	return count
}

func (u *UseCase) calculateUserStatistics(
	ctx context.Context,
) ([]UserAssignmentStats, error) {
	users, err := u.getAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	userStats := make([]UserAssignmentStats, 0, len(users))

	for _, user := range users {
		authoredPRs, err := u.repo.GetOpenPullRequestsByAuthor(ctx, user.ID)
		if err != nil {
			continue
		}

		assignedReviews, err := u.repo.GetUserAssignedPRs(ctx, user.ID)
		if err != nil {
			continue
		}

		var activeAssignments int64

		for _, review := range assignedReviews {
			pr, err := u.repo.GetPullRequestByID(ctx, review.PullRequestID)
			if err == nil && pr.Status == model.PullRequestStatusOpen {
				activeAssignments++
			}
		}

		var teamName *string

		teamMembers, err := u.repo.GetTeamMembersByTeamID(ctx, user.ID)

		if err == nil && len(teamMembers) > 0 {
			team, err := u.repo.GetTeamByID(ctx, teamMembers[0].TeamID)
			if err == nil {
				teamName = &team.Name
			}
		}

		userStats = append(userStats, UserAssignmentStats{
			UserID:             user.ExternalID,
			Username:           user.Username,
			TeamName:           teamName,
			TotalPRs:           int64(len(authoredPRs)),
			AssignedAsReviewer: int64(len(assignedReviews)),
			ActiveAssignments:  activeAssignments,
		})
	}

	return userStats, nil
}

func (u *UseCase) calculateTeamStatistics(
	ctx context.Context,
	prs []data.PullRequest,
) ([]TeamStatistics, error) {
	teams, err := u.getAllTeams(ctx)
	if err != nil {
		return nil, err
	}

	var teamStats []TeamStatistics

	for _, team := range teams {
		var teamPRs, openPRs int64

		for _, pr := range prs {
			author, err := u.repo.GetUserByID(ctx, pr.AuthorID)
			if err != nil {
				continue
			}

			isInTeam, err := u.isUserInTeam(ctx, author.ID, team.ID)

			if err == nil && isInTeam {
				teamPRs++

				if pr.Status == model.PullRequestStatusOpen {
					openPRs++
				}
			}
		}

		totalReviews, err := u.getTeamReviewAssignmentsCount(ctx, team.ID)
		if err != nil {
			totalReviews = 0
		}

		teamStats = append(teamStats, TeamStatistics{
			TeamName:     team.Name,
			TotalPRs:     teamPRs,
			OpenPRs:      openPRs,
			TotalReviews: totalReviews,
		})
	}

	return teamStats, nil
}

func (u *UseCase) calculateReviewerLoad(ctx context.Context) ([]ReviewerLoadStats, error) {
	users, err := u.getAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var reviewerLoad []ReviewerLoadStats

	for _, user := range users {
		assignedReviews, err := u.repo.GetUserAssignedPRs(ctx, user.ID)
		if err != nil {
			continue
		}

		var activeLoad int64

		for _, review := range assignedReviews {
			pr, err := u.repo.GetPullRequestByID(ctx, review.PullRequestID)
			if err == nil && pr.Status == model.PullRequestStatusOpen {
				activeLoad++
			}
		}

		if activeLoad > 0 {
			reviewerLoad = append(reviewerLoad, ReviewerLoadStats{
				UserID:   user.ExternalID,
				Username: user.Username,
				Load:     activeLoad,
			})
		}
	}

	return reviewerLoad, nil
}

func (u *UseCase) getAllUsers(ctx context.Context) ([]data.User, error) {
	var users []data.User

	teams, err := u.getAllTeams(ctx)
	if err != nil {
		return nil, err
	}

	for _, team := range teams {
		members, err := u.repo.GetTeamMembersByTeamID(ctx, team.ID)
		if err != nil {
			continue
		}

		for _, member := range members {
			user, err := u.repo.GetUserByID(ctx, member.UserID)
			if err == nil {
				users = append(users, user)
			}
		}
	}

	return users, nil
}

func (u *UseCase) getAllTeams(ctx context.Context) ([]data.Team, error) {
	teams, err := u.repo.GetAllTeams(ctx)
	if err != nil {
		return nil, err
	}

	return teams, nil
}

func (u *UseCase) isUserInTeam(ctx context.Context, userID, teamID uuid.UUID) (bool, error) {
	_, err := u.repo.GetTeamMemberByTeamAndUser(ctx, teamID, userID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *UseCase) getTeamReviewAssignmentsCount(
	ctx context.Context,
	teamID uuid.UUID,
) (int64, error) {
	teamMembers, err := u.repo.GetTeamMembersByTeamID(ctx, teamID)
	if err != nil {
		return 0, err
	}

	var total int64

	for _, member := range teamMembers {
		reviews, err := u.repo.GetUserAssignedPRs(ctx, member.UserID)
		if err == nil {
			total += int64(len(reviews))
		}
	}

	return total, nil
}
