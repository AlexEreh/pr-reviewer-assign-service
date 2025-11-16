package base

import "context"

type Tx interface {
	Executor
	Parent() Tx
	Commit(ctx context.Context) (err error)
	Rollback(ctx context.Context) (err error)
}

var _ Tx = (*Nop)(nil)

type Nop struct {
	Tx
}

func (n Nop) Parent() Tx {
	return n.Tx
}

func (n Nop) Commit(_ context.Context) (err error) {
	return nil
}

func (n Nop) Rollback(_ context.Context) (err error) {
	return nil
}
