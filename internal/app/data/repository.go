package data

import (
	"context"

	"github.com/google/uuid"
)

type TeamRepository interface {
	// GetTeamByID получает команду по ID.
	GetTeamByID(ctx context.Context, ID uuid.UUID) (Team, error)
	// GetTeamByExternalID получает команду по внешнему ID.
	GetTeamByExternalID(ctx context.Context, externalID string) (Team, error)
	// GetTeamByName получает команду по атрибуту "имя команды".
	GetTeamByName(ctx context.Context, teamName string) (Team, error)
	// GetAllTeams получает команду по атрибуту "имя команды".
	GetAllTeams(ctx context.Context) ([]Team, error)
	// CreateTeam создает новую команду.
	CreateTeam(ctx context.Context, team Team) (Team, error)
	// UpdateTeam обновляет данные команды.
	UpdateTeam(ctx context.Context, team Team) (Team, error)
}

type UserRepository interface {
	// GetUserByID получает пользователя по ID.
	GetUserByID(ctx context.Context, ID uuid.UUID) (User, error)
	// GetUserByExternalID получает пользователя по внешнему ID.
	GetUserByExternalID(ctx context.Context, externalID string) (User, error)
	// GetUserByName получает пользователя по атрибуту "никнейм".
	GetUserByName(ctx context.Context, userName string) (User, error)
	// GetUserByEmail получает пользователя по атрибуту "электронная почта".
	GetUserByEmail(ctx context.Context, email string) (User, error)
	// CreateUser создает нового пользователя.
	CreateUser(ctx context.Context, user User) (User, error)
	// UpdateUser обновляет данные пользователя.
	UpdateUser(ctx context.Context, user User) (User, error)
	// SetUserActive устанавливает флаг активности пользователя.
	SetUserActive(ctx context.Context, userID uuid.UUID, isActive bool) (User, error)
}

type TeamMemberRepository interface {
	// GetTeamMemberByID получает связь пользователя с командой по ID.
	GetTeamMemberByID(ctx context.Context, ID uuid.UUID) (TeamMember, error)
	// GetTeamMemberByTeamAndUser получает связь пользователя с командой по команде и пользователю.
	GetTeamMemberByTeamAndUser(ctx context.Context, teamID, userID uuid.UUID) (TeamMember, error)
	// CreateTeamMember создает новую связь пользователя с командой.
	CreateTeamMember(ctx context.Context, teamMember TeamMember) (TeamMember, error)
	// UpdateTeamMemberRole обновляет роль участника команды.
	UpdateTeamMemberRole(
		ctx context.Context,
		teamMemberID uuid.UUID,
		role string,
	) (TeamMember, error)
	// DeleteTeamMember удаляет участника из команды.
	DeleteTeamMember(ctx context.Context, teamMemberID uuid.UUID) error
	// GetTeamMembersByTeamID возвращает всех участников команды.
	GetTeamMembersByTeamID(ctx context.Context, teamID uuid.UUID) ([]TeamMember, error)
	// GetTeamMembersByUserID возвращает все участия в командах для пользователя.
	GetTeamMembersByUserID(ctx context.Context, userID uuid.UUID) ([]TeamMember, error)
}

type PullRequestRepository interface {
	// GetPullRequestByID получает PR по ID.
	GetPullRequestByID(ctx context.Context, ID uuid.UUID) (PullRequest, error)
	// GetPullRequestByExternalID получает PR по внешнему ID.
	GetPullRequestByExternalID(ctx context.Context, externalID string) (PullRequest, error)
	// CreatePullRequest создает новый PR.
	CreatePullRequest(ctx context.Context, pr PullRequest) (PullRequest, error)
	// UpdatePullRequest обновляет данные PR.
	UpdatePullRequest(ctx context.Context, pr PullRequest) (PullRequest, error)
	// MergePullRequest помечает PR как мерженный.
	MergePullRequest(ctx context.Context, prID uuid.UUID) (PullRequest, error)
	// GetOpenPullRequestsByAuthor возвращает открытые PR автора.
	GetOpenPullRequestsByAuthor(ctx context.Context, authorID uuid.UUID) ([]PullRequest, error)
	// GetPullRequestsByStatus возвращает PR автора.
	GetPullRequestsByStatus(ctx context.Context, status string) ([]PullRequest, error)
}

type PRReviewerRepository interface {
	// GetPRReviewerByID получает назначение ревьювера по ID.
	GetPRReviewerByID(ctx context.Context, ID uuid.UUID) (PRReviewer, error)
	// CreatePRReviewer создает новое назначение ревьювера.
	CreatePRReviewer(ctx context.Context, reviewer PRReviewer) (PRReviewer, error)
	// GetCurrentReviewers возвращает текущих ревьюверов PR.
	GetCurrentReviewers(ctx context.Context, prID uuid.UUID) ([]PRReviewer, error)
	// GetUserAssignedPRs возвращает PR, на которые назначен пользователь как ревьювер.
	GetUserAssignedPRs(ctx context.Context, userID uuid.UUID) ([]PRReviewer, error)
	// UpdatePRReviewer обновляет данные назначения ревьювера
	UpdatePRReviewer(ctx context.Context, reviewer PRReviewer) (PRReviewer, error)
}

type PRReviewerHistoryRepository interface {
	// GetPRReviewerHistoryByID получает запись истории по ID.
	GetPRReviewerHistoryByID(ctx context.Context, ID uuid.UUID) (PRReviewerHistory, error)
	// CreatePRReviewerHistory создает новую запись истории изменений ревьюверов.
	CreatePRReviewerHistory(
		ctx context.Context,
		history PRReviewerHistory,
	) (PRReviewerHistory, error)
	// GetPRReviewerHistory возвращает историю изменений для PR.
	GetPRReviewerHistory(ctx context.Context, prID uuid.UUID) ([]PRReviewerHistory, error)
}

// Repository объединяет все репозитории для удобства использования
type Repository interface {
	TeamRepository
	UserRepository
	TeamMemberRepository
	PullRequestRepository
	PRReviewerRepository
	PRReviewerHistoryRepository
}
