package data

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type (
	UserInternalID              = uuid.UUID
	TeamInternalID              = uuid.UUID
	PullRequestInternalID       = uuid.UUID
	PRReviewerHistoryInternalID = uuid.UUID
	PRReviewerInternalID        = uuid.UUID
)

type (
	UserExternalID        = string
	TeamExternalID        = string
	PullRequestExternalID = string
)

type User struct {
	ID         UserInternalID
	ExternalID UserExternalID
	Username   string
	Email      string
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Team struct {
	ID          TeamInternalID
	ExternalID  TeamExternalID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TeamMember struct {
	ID        uuid.UUID
	TeamID    TeamInternalID
	UserID    UserInternalID
	Role      string
	CreatedAt time.Time
}

type PullRequest struct {
	ID                PullRequestInternalID
	ExternalID        PullRequestExternalID
	Title             string
	Description       string
	AuthorID          uuid.UUID
	Status            string
	NeedMoreReviewers bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
	MergedAt          sql.NullTime
}

type PRReviewer struct {
	ID            PRReviewerInternalID
	PullRequestID PullRequestInternalID
	ReviewerID    UserInternalID
	TeamID        TeamInternalID
	AssignedAt    time.Time
	ReplacedAt    sql.NullTime
	IsCurrent     bool
}

type PRReviewerHistory struct {
	ID            PRReviewerHistoryInternalID
	PullRequestID PullRequestInternalID
	OldReviewerID sql.Null[UserInternalID]
	NewReviewerID UserInternalID
	ChangedBy     sql.Null[UserInternalID]
	ChangedAt     time.Time
	Reason        string
}
