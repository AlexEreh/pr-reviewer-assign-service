package postgres

import (
	"pr-reviewer-assign-service/internal/app/data"
	"pr-reviewer-assign-service/pkg/txman"
)

type Repository struct {
	TeamRepository
	UserRepository
	TeamMemberRepository
	PullRequestRepository
	PRReviewerRepository
	PRReviewerHistoryRepository
}

func NewRepository(txMan txman.Manager) *Repository {
	return &Repository{
		TeamRepository:              TeamRepository{txMan: txMan},
		UserRepository:              UserRepository{txMan: txMan},
		TeamMemberRepository:        TeamMemberRepository{txMan: txMan},
		PullRequestRepository:       PullRequestRepository{txMan: txMan},
		PRReviewerRepository:        PRReviewerRepository{txMan: txMan},
		PRReviewerHistoryRepository: PRReviewerHistoryRepository{txMan: txMan},
	}
}

var _ data.Repository = (*Repository)(nil)
