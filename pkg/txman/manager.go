package txman

import (
	"context"
	"database/sql"
	"fmt"

	"pr-reviewer-assign-service/pkg/txman/base"
)

type Manager struct {
	base base.Manager
}

func New(db *sql.DB, options ...Option) Manager {
	return Manager{base.New(
		fmt.Sprintf("%p", db),
		&baseDb{db},
		options...,
	)}
}

func (m Manager) Transactional(
	ctx context.Context,
	f func(ctx context.Context) error,
	options ...Option,
) (err error) {
	return m.base.Transactional(ctx, f, options...)
}

func (m Manager) Executor(ctx context.Context) (executor Executor) {
	return m.base.Executor(ctx).(Executor)
}

type BeginError = base.BeginError

type CommitError = base.CommitError

type RollbackError = base.RollbackError
