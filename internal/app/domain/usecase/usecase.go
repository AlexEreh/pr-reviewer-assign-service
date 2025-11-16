package usecase

import (
	"pr-reviewer-assign-service/internal/app/data"
	"pr-reviewer-assign-service/pkg/txman"
)

type UseCase struct {
	repo  data.Repository
	txMan txman.Manager
}

func New(repo data.Repository, txMan txman.Manager) *UseCase {
	return &UseCase{
		repo:  repo,
		txMan: txMan,
	}
}
