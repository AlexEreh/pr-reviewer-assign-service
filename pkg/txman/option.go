package txman

import (
	"database/sql"

	"pr-reviewer-assign-service/pkg/txman/base"
)

type Option = base.Option

const propagationKey = "propagation"

func WithPropagation(propagation Propagation) Option {
	return Option{
		Key:   propagationKey,
		Value: propagation,
	}
}

func getPropagation(options base.Options) (propagation Propagation) {
	propagation, _ = options.Value(propagationKey).(Propagation)
	return propagation
}

type Propagation int

const (
	Never Propagation = iota - 2
	Supports
	Required
	Nested
	Mandatory
)

const isolationLevelKey = "isolation"

func WithIsolationLevel(level sql.IsolationLevel) (option Option) {
	return base.Option{
		Key:   isolationLevelKey,
		Value: level,
	}
}

func getIsolationLevel(options base.Options) (level sql.IsolationLevel) {
	level, _ = options.Value(isolationLevelKey).(sql.IsolationLevel)
	return level
}

const readonlyKey = "readonly"

func WithReadonly(readonly bool) (option Option) {
	return base.Option{
		Key:   readonlyKey,
		Value: readonly,
	}
}

func getReadonly(options base.Options) (readonly bool) {
	readonly, _ = options.Value(readonlyKey).(bool)
	return readonly
}
