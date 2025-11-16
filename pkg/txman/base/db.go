package base

import "context"

type Db interface {
	Executor
	Begin(ctx context.Context, tx Tx, options Options) (newTx Tx, err error)
}
