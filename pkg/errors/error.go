package errors

import "maps"

type Code string

type Error struct {
	code       Code
	message    string
	cause      error
	paramsMap  map[string]any
	stackTrace StackTrace
}

const skipCallPointStackFramesCount = 2

var _ error = (*Error)(nil)

func New(template Template, params ...Param) error {
	return newError(template, nil, params)
}

func Wrap(err error, template Template, params ...Param) error {
	return newError(template, err, params)
}

func Is(err error, template Template) bool {
	if err == nil {
		return false
	}

	e, ok := err.(*Error)
	if !ok {
		return false
	}

	return e.Code() == template.Code
}

func newError(template Template, cause error, params Params) *Error {
	return &Error{
		code:       template.Code,
		message:    template.Message,
		cause:      cause,
		paramsMap:  mergeParamMaps(template.Params.toMap(), params.toMap()),
		stackTrace: trace(skipCallPointStackFramesCount),
	}
}

func mergeParamMaps(paramMaps ...map[string]any) map[string]any {
	size := 0
	for _, p := range paramMaps {
		size += len(p)
	}

	mergedParams := make(map[string]any, size)

	for _, m := range paramMaps {
		maps.Copy(mergedParams, m)
	}

	return mergedParams
}

func (e *Error) Code() Code {
	return e.code
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Cause() error {
	return e.cause
}

func (e *Error) Get(key string) any {
	if e.paramsMap == nil {
		return nil
	}

	return e.paramsMap[key]
}

func (e *Error) Params() Params {
	params := make(Params, 0, len(e.paramsMap))
	for key, value := range e.paramsMap {
		params = append(params, Param{key, value})
	}

	return params
}

func (e *Error) StackTrace() StackTrace {
	return e.stackTrace
}
