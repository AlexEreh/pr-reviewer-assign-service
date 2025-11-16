package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"pr-reviewer-assign-service/internal/app/delivery/http/api"
	"pr-reviewer-assign-service/internal/app/delivery/http/api/health"
	"pr-reviewer-assign-service/internal/app/delivery/http/api/pullrequests"
	"pr-reviewer-assign-service/internal/app/delivery/http/api/statistics"
	"pr-reviewer-assign-service/internal/app/delivery/http/api/teams"
	"pr-reviewer-assign-service/internal/app/delivery/http/api/users"
)

type E2ETestSuite struct {
	suite.Suite

	client    *http.Client
	apiClient api.Client
}

const (
	baseURL       = "http://localhost:8080"
	clientTimeout = 30 * time.Second
)

func TestE2ESuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}

func (s *E2ETestSuite) SetupSuite() {
	s.client = &http.Client{ //nolint:exhaustruct
		Timeout: clientTimeout,
	}

	s.apiClient = api.NewClient(s.client, baseURL)

	s.waitForService()
}

func (s *E2ETestSuite) waitForService() {
	for i := 0; i < 2; i++ {
		_, err := s.apiClient.Health().ReadyZ(s.T().Context(), health.ReadyZParams{})
		if err == nil {
			fmt.Println("Service is ready!")

			return
		} else {
			fmt.Println("Error: ", err)
		}

		time.Sleep(1 * time.Second)
	}

	s.T().Fatal("Service did not become ready in time")
}

// TestHealth проверяет health endpoint.
func (s *E2ETestSuite) TestHealth() {
	result, err := s.apiClient.Health().ReadyZ(s.T().Context(), health.ReadyZParams{})

	s.NoError(err)
	s.True(result.OK)
}

// TestTeamLifecycle тестирует создание и получение команды.
func (s *E2ETestSuite) TestTeamLifecycle() {
	teamName := fmt.Sprintf("team-%d", time.Now().UnixNano())

	firstUserID := fmt.Sprintf("user-1-%d", time.Now().UnixNano())
	firstUserName := fmt.Sprintf("Test User 1 %d", time.Now().UnixNano())
	secondUserID := fmt.Sprintf("user-2-%d", time.Now().UnixNano())
	secondUserName := fmt.Sprintf("Test User 2 %d", time.Now().UnixNano())

	result, err := s.apiClient.Teams().AddTeam(s.T().Context(), teams.AddTeamParams{
		TeamName: teamName,
		Members: []teams.AddTeamParamsUser{
			{
				UserID:   firstUserID,
				UserName: firstUserName,
				IsActive: true,
			},
			{
				UserID:   secondUserID,
				UserName: secondUserName,
				IsActive: true,
			},
		},
	})

	s.NoError(err)
	s.Equal(teamName, result.Team.TeamName)

	getTeamResult, err := s.apiClient.Teams().
		GetTeam(s.T().Context(), teams.GetTeamParams{TeamName: teamName})

	s.NoError(err)
	s.Equal(teamName, getTeamResult.TeamName)
	s.Len(getTeamResult.Members, 2)
}

// TestUserActivation тестирует активацию/деактивацию пользователя.
func (s *E2ETestSuite) TestUserActivation() {
	teamName := fmt.Sprintf("team-activation-%d", time.Now().UnixNano())

	userID := fmt.Sprintf("user-activation-%d", time.Now().UnixNano())
	userName := fmt.Sprintf("Test Activation User %d", time.Now().UnixNano())

	_, err := s.apiClient.Teams().AddTeam(s.T().Context(), teams.AddTeamParams{
		TeamName: teamName,
		Members: []teams.AddTeamParamsUser{
			{
				UserID:   userID,
				UserName: userName,
				IsActive: true,
			},
		},
	})

	s.NoError(err)

	deactivateResult, err := s.apiClient.Users().
		SetIsActive(s.T().Context(), users.SetIsActiveParams{
			UserID:   userID,
			IsActive: false,
		})

	s.NoError(err)
	s.Equal(userID, deactivateResult.User.UserID)
	s.Equal(false, deactivateResult.User.IsActive)
}

// TestPullRequestLifecycle тестирует полный жизненный цикл PR.
func (s *E2ETestSuite) TestPullRequestLifecycle() {
	teamName := fmt.Sprintf("team-pr-%d", time.Now().UnixNano())
	authorID := fmt.Sprintf("author-pr-test-%d", time.Now().UnixNano())
	authorName := fmt.Sprintf("PR Author %d", time.Now().UnixNano())
	reviewerID := fmt.Sprintf("reviewer-pr-test-%d", time.Now().UnixNano())
	reviewerName := fmt.Sprintf("PR Reviewer %d", time.Now().UnixNano())

	_, err := s.apiClient.Teams().AddTeam(s.T().Context(), teams.AddTeamParams{
		TeamName: teamName,
		Members: []teams.AddTeamParamsUser{
			{
				UserID:   authorID,
				UserName: authorName,
				IsActive: true,
			},
			{
				UserID:   reviewerID,
				UserName: reviewerName,
				IsActive: true,
			},
		},
	})
	s.NoError(err)

	prID := fmt.Sprintf("pr-%d", time.Now().UnixNano())
	prName := fmt.Sprintf("Test PR for E2E %d", time.Now().UnixNano())
	createResult, err := s.apiClient.PR().CreatePR(s.T().Context(), pullrequests.CreatePRParams{
		PullRequestID:   prID,
		PullRequestName: prName,
		AuthorID:        authorID,
	})
	s.NoError(err)

	s.Equal(prID, createResult.PR.PullRequestID)
	s.Equal("OPEN", createResult.PR.Status)
	s.NotEmpty(createResult.PR.AssignedReviewers)

	reviewResult, err := s.apiClient.Users().GetReviewPRs(s.T().Context(), users.GetReviewPRsParams{
		UserID: reviewerID,
	})
	s.NoError(err)
	s.Equal(reviewerID, reviewResult.UserID)
	s.Len(reviewResult.PullRequests, 1)

	mergeResult, err := s.apiClient.PR().MergePR(s.T().Context(), pullrequests.MergePRParams{
		PullRequestID: prID,
	})
	s.NoError(err)

	s.Equal("MERGED", mergeResult.PR.Status)
	s.NotNil(mergeResult.PR.MergedAt)
}

// TestReviewerReassignment тестирует переназначение ревьювера
func (s *E2ETestSuite) TestReviewerReassignment() {
	teamName := fmt.Sprintf("team-reassign-%d", time.Now().UnixNano())

	authorID := fmt.Sprintf("author-reassign-test-%d", time.Now().UnixNano())
	authorName := fmt.Sprintf("Reassign Author %d", time.Now().UnixNano())

	oldReviewerID := fmt.Sprintf("old-reviewer-test-%d", time.Now().UnixNano())
	oldReviewerName := fmt.Sprintf("Old Reviewer %d", time.Now().UnixNano())

	newReviewerID := fmt.Sprintf("new-reviewer-test-%d", time.Now().UnixNano())
	newReviewerName := fmt.Sprintf("New Reviewer %d", time.Now().UnixNano())

	_, err := s.apiClient.Teams().AddTeam(s.T().Context(), teams.AddTeamParams{
		TeamName: teamName,
		Members: []teams.AddTeamParamsUser{
			{
				UserID:   authorID,
				UserName: authorName,
				IsActive: true,
			},
			{
				UserID:   oldReviewerID,
				UserName: oldReviewerName,
				IsActive: true,
			},
			{
				UserID:   newReviewerID,
				UserName: newReviewerName,
				IsActive: true,
			},
		},
	})
	s.NoError(err)

	prID := fmt.Sprintf("pr-reassign-%d", time.Now().UnixNano())
	prName := fmt.Sprintf("Reassignment Test PR %d", time.Now().UnixNano())
	_, err = s.apiClient.PR().CreatePR(s.T().Context(), pullrequests.CreatePRParams{
		PullRequestID:   prID,
		PullRequestName: prName,
		AuthorID:        authorID,
	})
	s.NoError(err)

	reassignResult, err := s.apiClient.PR().
		ReassignPR(s.T().Context(), pullrequests.ReassignPRParams{
			PullRequestID: prID,
			OldReviewerID: oldReviewerID,
		})
	s.NoError(err)

	s.NotEmpty(reassignResult.ReplacedBy)
	s.NotEqual(oldReviewerID, reassignResult.ReplacedBy)
}

// TestStatistics тестирует эндпоинт статистики
func (s *E2ETestSuite) TestStatistics() {
	getStatsResult, err := s.apiClient.Statistics().
		GetStatistics(s.T().Context(), statistics.GetStatisticsParams{})
	s.NoError(err)

	s.Equal(int64(0), getStatsResult.OpenPRs)
	s.Equal(int64(0), getStatsResult.MergedPRs)
	s.Equal(0, len(getStatsResult.UserAssignments))
	s.Equal(0, len(getStatsResult.TeamStats))
	s.Equal(0, len(getStatsResult.ReviewerLoad))
}

// TestErrorScenarios тестирует обработку ошибок
func (s *E2ETestSuite) TestErrorScenarios() {
	teamName := fmt.Sprintf("duplicate-team-test-%d", time.Now().UnixNano())
	userID := fmt.Sprintf("user-dup-%d", time.Now().UnixNano())
	userName := fmt.Sprintf("Duplicate User %d", time.Now().UnixNano())
	_, err := s.apiClient.Teams().AddTeam(s.T().Context(), teams.AddTeamParams{
		TeamName: teamName,
		Members: []teams.AddTeamParamsUser{
			{
				UserID:   userID,
				UserName: userName,
				IsActive: true,
			},
		},
	})
	s.NoError(err)

	_, err = s.apiClient.Teams().AddTeam(s.T().Context(), teams.AddTeamParams{
		TeamName: teamName,
		Members: []teams.AddTeamParamsUser{
			{
				UserID:   userID + " ",
				UserName: userName + " ",
				IsActive: true,
			},
		},
	})
	s.Error(err)
	s.Contains(err.Error(), "TEAM_EXISTS")

	// Попытка получить несуществующую команду
	_, err = s.apiClient.Teams().GetTeam(s.T().Context(), teams.GetTeamParams{
		TeamName: "nonexistent-team",
	})
	s.Error(err)
	s.Contains(err.Error(), "NOT_FOUND")

	// Попытка создать PR с несуществующим автором
	_, err = s.apiClient.PR().CreatePR(s.T().Context(), pullrequests.CreatePRParams{
		PullRequestID:   fmt.Sprintf("pr-error-test-%d", time.Now().UnixNano()),
		PullRequestName: fmt.Sprintf("Error Test PR %d", time.Now().UnixNano()),
		AuthorID:        "nonexistent-author",
	})
	s.Error(err)
	s.Contains(err.Error(), "NOT_FOUND")
}

// TestIdempotentMerge тестирует идемпотентность мержа PR
func (s *E2ETestSuite) TestIdempotentMerge() {
	teamName := fmt.Sprintf("team-idempotent-%d", time.Now().UnixNano())
	authorId := fmt.Sprintf("author-idempotent-test-%d", time.Now().UnixNano())
	authorName := fmt.Sprintf("Idempotent Author %d", time.Now().UnixNano())

	_, err := s.apiClient.Teams().AddTeam(s.T().Context(), teams.AddTeamParams{
		TeamName: teamName,
		Members: []teams.AddTeamParamsUser{
			{
				UserID:   authorId,
				UserName: authorName,
				IsActive: true,
			},
		},
	})

	s.NoError(err)

	prID := uuid.NewString()
	prName := uuid.NewString()

	_, err = s.apiClient.PR().CreatePR(s.T().Context(), pullrequests.CreatePRParams{
		PullRequestID:   prID,
		PullRequestName: prName,
		AuthorID:        authorId,
	})

	s.NoError(err)

	firstResult, err := s.apiClient.PR().MergePR(s.T().Context(), pullrequests.MergePRParams{
		PullRequestID: prID,
	})

	s.NoError(err)

	secondResult, err := s.apiClient.PR().MergePR(s.T().Context(), pullrequests.MergePRParams{
		PullRequestID: prID,
	})

	s.NoError(err)

	s.Equal(firstResult.PR.MergedAt, secondResult.PR.MergedAt)
}
