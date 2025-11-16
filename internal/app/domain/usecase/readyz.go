package usecase

import (
	"context"

	"pr-reviewer-assign-service/pkg/errors"
)

type ReadyZParams struct{}

type ReadyZResult struct {
	OK bool
}

func (u *UseCase) ReadyZ(ctx context.Context, _ ReadyZParams) (ReadyZResult, error) {
	var i int

	err := u.txMan.Executor(ctx).QueryRowContext(ctx, "select 1;").Scan(&i)
	if err != nil {
		return ReadyZResult{OK: false}, nil
	}

	if i != 1 {
		return ReadyZResult{OK: false}, errors.New(errors.InternalError)
	}

	return ReadyZResult{OK: true}, nil
}
