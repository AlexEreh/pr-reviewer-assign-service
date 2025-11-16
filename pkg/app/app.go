package app

import (
	"context"
	"errors"
	"slices"
	"sync"
	"time"
)

type App struct {
	lock                   sync.Mutex
	ctx                    context.Context //nolint
	done                   chan struct{}
	errs                   []error
	isInitialized          bool
	afterInitHooks         []func() (err error)
	isWaitingForCompletion bool
	goroutinesCount        int
	goroutines             chan struct{}
	isCompleted            bool
	completeHooks          []func() (err error)
}

var _ context.Context = (*App)(nil)

var ErrCompleted = errors.New("completed")

func (a *App) Deadline() (deadline time.Time, ok bool) {
	return a.ctx.Deadline()
}

func (a *App) Done() (ch <-chan struct{}) {
	return a.done
}

func (a *App) Err() (err error) {
	a.lock.Lock()

	if a.isCompleted {
		if len(a.errs) > 0 {
			err = a.errs[0]
		} else {
			err = ErrCompleted
		}
	}

	a.lock.Unlock()

	return err
}

func (a *App) Value(key any) any {
	return a.ctx.Value(key)
}

func Run(ctx context.Context, init func(app *App) (err error)) (errs []error) {
	app := &App{
		lock:                   sync.Mutex{},
		ctx:                    ctx,
		done:                   make(chan struct{}),
		errs:                   nil,
		isInitialized:          false,
		isWaitingForCompletion: false,
		goroutinesCount:        0,
		goroutines:             make(chan struct{}),
		isCompleted:            false,
		completeHooks:          nil,
		afterInitHooks:         nil,
	}

	return app.run(init)
}

func (a *App) OnComplete(hook func() error) {
	a.lock.Lock()

	if a.isWaitingForCompletion {
		a.lock.Unlock()
		panic("must be called inside the Init() function or inside the AfterInit() hook")
	}

	a.completeHooks = append(a.completeHooks, hook)
	a.lock.Unlock()
}

func (a *App) Go(f func() error) {
	a.lock.Lock()

	if a.isWaitingForCompletion {
		a.lock.Unlock()
		panic("must be called inside the Init() function or inside the AfterInit() hook")
	}

	a.goroutinesCount++
	a.lock.Unlock()

	go func() {
		defer func() {
			a.lock.Lock()
			a.goroutinesCount--

			if a.goroutinesCount == 0 && a.isWaitingForCompletion {
				close(a.goroutines)
			}

			a.lock.Unlock()
		}()

		err := f()
		if err == nil {
			return
		}

		a.lock.Lock()
		a.addError(err)
		a.complete()
		a.lock.Unlock()
	}()
}

func (a *App) AfterInit(hook func() error) {
	a.lock.Lock()

	if a.isInitialized {
		a.lock.Unlock()
		panic("must be called inside the Init() function")
	}

	a.afterInitHooks = append(a.afterInitHooks, hook)
	a.lock.Unlock()
}

func (a *App) run(init func(app *App) (err error)) []error {
	go a.notifyCompletion()

	a.init(init)
	a.invokeAfterInitHooks()
	a.waitForCompletion()
	a.invokeCompleteHooks()
	a.waitGoroutines()

	return a.errors()
}

func (a *App) notifyCompletion() {
	select {
	case <-a.ctx.Done():
		a.lock.Lock()
		a.addError(a.ctx.Err())
		a.complete()
		a.lock.Unlock()
	case <-a.goroutines:
		a.lock.Lock()
		a.complete()
		a.lock.Unlock()
	}
}

func (a *App) init(init func(app *App) (err error)) {
	err := init(a)
	a.lock.Lock()
	a.isInitialized = true

	if err != nil {
		a.addError(err)
		a.complete()
	}

	a.lock.Unlock()
}

func (a *App) invokeAfterInitHooks() {
	a.lock.Lock()

	hooks := a.afterInitHooks

	a.afterInitHooks = nil

	a.lock.Unlock()

	for hookIdx := range hooks {
		a.lock.Lock()
		isCompleted := a.isCompleted
		a.lock.Unlock()

		if isCompleted {
			break
		}

		err := hooks[hookIdx]()
		if err == nil {
			continue
		}

		a.lock.Lock()
		a.addError(err)
		a.complete()
		a.lock.Unlock()
	}
}

func (a *App) waitForCompletion() {
	a.lock.Lock()
	a.isWaitingForCompletion = true

	if a.goroutinesCount == 0 {
		close(a.goroutines)
	}

	a.lock.Unlock()
	<-a.done
}

func (a *App) invokeCompleteHooks() {
	a.lock.Lock()

	completeHooks := a.completeHooks

	a.completeHooks = nil

	a.lock.Unlock()

	for i := len(completeHooks) - 1; i >= 0; i-- {
		err := completeHooks[i]()
		if err == nil {
			continue
		}

		a.lock.Lock()
		a.addError(err)
		a.lock.Unlock()
	}
}

func (a *App) waitGoroutines() {
	<-a.goroutines
}

func (a *App) complete() {
	if a.isCompleted {
		return
	}

	a.isCompleted = true

	close(a.done)
}

func (a *App) addError(err error) {
	a.errs = append(a.errs, err)
}

func (a *App) errors() []error {
	a.lock.Lock()
	errs := slices.Clone(a.errs)
	a.lock.Unlock()

	return errs
}
