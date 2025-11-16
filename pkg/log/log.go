package log

import (
	"context"
)

type ctxLogger struct{}

var defaultLogger Logger

func LoggerFromCtx(ctx context.Context) Logger {
	l, ok := ctx.Value(ctxLogger{}).(Logger)
	if ok {
		return l
	}

	return defaultLogger
}

func WithLogger(ctx context.Context, lg Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, lg)
}

func SetLogger(lg Logger, f func(key, value any)) {
	f(ctxLogger{}, lg)
}
