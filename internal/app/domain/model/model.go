package model

type TeamMember struct {
	UserID   string
	Username string
	IsActive bool
}

type Team struct {
	TeamName string
	Members  []TeamMember
}

type User struct {
	UserID   string
	UserName string
	IsActive bool
	TeamName string
}

type PullRequest struct {
	PullRequestShort

	AssignedReviewers []string
}

type PullRequestShort struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
	Status          string
}

type PullRequestStatus = string

const (
	PullRequestStatusOpen   = "OPEN"
	PullRequestStatusMerged = "MERGED"
)

type PRReviewerHistoryChangeReason = string

const (
	PRReviewerHistoryChangeReasonInitial      = "initial"
	PRReviewerHistoryChangeReasonReassignment = "reassignment"
)
