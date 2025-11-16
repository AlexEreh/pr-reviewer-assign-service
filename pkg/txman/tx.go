package txman

import (
	"context"
	"database/sql"

	"pr-reviewer-assign-service/pkg/txman/base"
)

type baseTx struct {
	tx      *sql.Tx
	options *sql.TxOptions
}

func (t baseTx) Executor() any {
	return t.tx
}

func (t baseTx) Parent() base.Tx {
	return nil
}

func (t baseTx) Commit(_ context.Context) error {
	return t.tx.Commit() //nolint:wrapcheck // Пока не требуется
}

func (t baseTx) Rollback(_ context.Context) error {
	return t.tx.Rollback() //nolint:wrapcheck // Пока не требуется
}

func nop(tx base.Tx) base.Nop {
	return base.Nop{Tx: tx}
}
