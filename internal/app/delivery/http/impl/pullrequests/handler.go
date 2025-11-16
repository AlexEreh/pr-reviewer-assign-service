package pullrequests

import "pr-reviewer-assign-service/internal/app/domain/usecase"

type Handler struct {
	useCase *usecase.UseCase
}

func NewHandler(useCase *usecase.UseCase) *Handler {
	return &Handler{useCase: useCase}
}
