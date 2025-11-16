package base

import (
	"context"
	"fmt"
)

type Manager struct {
	key     any
	db      Db
	options []Option
}

func New(key any, db Db, options ...Option) (manager Manager) {
	return Manager{
		key:     key,
		db:      db,
		options: options,
	}
}

func (m Manager) Transactional(
	ctx context.Context,
	f func(ctx context.Context) (err error),
	options ...Option,
) (err error) {
	tx := m.extractTxFromContext(ctx)

	options = extend(nil, m.options, options)

	tx, err = m.db.Begin(ctx, tx, options)
	if err != nil {
		return BeginError{err}
	}

	ctx = m.putTxToContext(ctx, tx)
	if tx == nil {
		return f(ctx)
	}

	return m.transactional(ctx, tx, f)
}

func (m Manager) transactional(
	ctx context.Context,
	tx Tx,
	f func(ctx context.Context) (err error),
) (err error) {
	panicked := true

	defer func() {
		if panicked {
			_ = tx.Rollback(ctx)
		}
	}()

	err = f(ctx)
	panicked = false

	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			return RollbackError{err, rollbackErr}
		}

		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return CommitError{err}
	}

	return nil
}

func (m Manager) Executor(ctx context.Context) (executor any) {
	tx := m.extractTxFromContext(ctx)
	if tx != nil {
		return tx.Executor()
	}

	return m.db.Executor()
}

func (m Manager) extractTxFromContext(ctx context.Context) (tx Tx) {
	contextAny := ctx.Value(m.key)
	if contextAny == nil {
		return nil
	}

	tx, _ = contextAny.(Tx)

	return tx
}

func (m Manager) putTxToContext(ctx context.Context, tx Tx) (newCtx context.Context) {
	return context.WithValue(ctx, m.key, tx)
}

type BeginError struct {
	err error
}

func (e BeginError) Cause() (err error) {
	return e.err
}

func (e BeginError) Unwrap() (err error) {
	return e.err
}

func (e BeginError) Error() (err string) {
	return fmt.Sprintf("begin: %v", e.err)
}

type CommitError struct {
	err error
}

func (e CommitError) Cause() (err error) {
	return e.err
}

func (e CommitError) Unwrap() (err error) {
	return e.err
}

func (e CommitError) Error() (err string) {
	return fmt.Sprintf("commit: %v", e.err)
}

type RollbackError struct {
	txErr       error
	rollbackErr error
}

func (e RollbackError) Tx() (err error) {
	return e.txErr
}

func (e RollbackError) Rollback() (err error) {
	return e.rollbackErr
}

func (e RollbackError) Cause() (err error) {
	return e.rollbackErr
}

func (e RollbackError) Unwrap() (err error) {
	return e.rollbackErr
}

func (e RollbackError) Errors() (errs []error) {
	return []error{e.txErr, e.rollbackErr}
}

func (e RollbackError) Error() (err string) {
	return fmt.Sprintf("tx: %v; rollback: %v", e.txErr, e.rollbackErr)
}
