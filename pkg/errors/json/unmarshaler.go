package json

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"maps"
	"unicode"
	"unicode/utf8"

	stderrors "errors"

	"pr-reviewer-assign-service/pkg/errors"
)

type UnmarshalerOption func(d *UnmarshalerConfig)

func Unmarshal(bytes []byte, options ...UnmarshalerOption) (err error) {
	e := Error{
		Error:       err,
		Unmarshaler: options,
	}
	merr := xml.Unmarshal(bytes, &e)
	if merr == nil {
		return e.Error
	}
	return stderrors.New(string(bytes))
}

var _ json.Unmarshaler = (*Error)(nil)

func (e *Error) UnmarshalJSON(bytes []byte) (err error) {
	if isNilError(bytes) {
		e.Error = nil
		return nil
	}
	e.Error, err = unmarshalConcreteError(bytes, e.Unmarshaler...)
	if err == nil {
		return nil
	}
	e.Error = unmarshalAbstractError(bytes)
	return nil
}

func unmarshalAbstractError(bytes []byte) (e error) {
	var str string
	err := json.Unmarshal(bytes, &str)
	if err == nil {
		return stderrors.New(str)
	}
	return stderrors.New(string(bytes))
}

func unmarshalConcreteError(bytes []byte, options ...UnmarshalerOption) (e, err error) {
	cfg := applyUnmarshalerOptions(options)
	data := map[string]json.RawMessage{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	var code errors.Code
	var codeFound bool
	var message string
	var messageFound bool
	var cause error
	var params errors.Params
	for jsonKey, jsonValue := range data {
		var ok bool
		key := cfg.UnmarshalKey(jsonKey)
		switch key {
		case keyCode:
			codeFound = true
			codeValue, err := unmarshalParam(cfg, keyCode, jsonValue)
			if err != nil {
				return nil, keyUnmarshalError{keyCode, err}
			}
			code, ok = codeValue.(errors.Code)
			if !ok {
				return nil, keyCastError{keyCode}
			}
		case keyMessage:
			messageFound = true
			messageValue, err := unmarshalParam(cfg, keyMessage, jsonValue)
			if err != nil {
				return nil, keyUnmarshalError{keyMessage, err}
			}
			message, ok = messageValue.(string)
			if !ok {
				return nil, keyCastError{keyMessage}
			}
		case keyCause:
			causeValue, err := unmarshalParam(cfg, keyCause, jsonValue)
			if err != nil {
				return nil, keyUnmarshalError{keyCause, err}
			}
			cause, ok = causeValue.(error)
			if !ok {
				return nil, keyCastError{keyCause}
			}
		case keyStackTrace:
		default:
			value, err := unmarshalParam(cfg, key, jsonValue)
			if err != nil {
				return nil, keyUnmarshalError{key, err}
			}
			params = append(params, errors.Param{Name: key, Value: value})
		}
	}
	if !codeFound {
		return nil, keyMissingError{keyCode}
	}
	if !messageFound {
		return nil, keyMissingError{keyMessage}
	}
	return errors.Wrap(cause, errors.Template{
		Code:    code,
		Message: message,
	}, params...), nil
}

func isNilError(bytes []byte) bool {
	return len(bytes) == 0 || string(bytes) == "null"
}

func applyUnmarshalerOptions(options []UnmarshalerOption) UnmarshalerConfig {
	d := UnmarshalerConfig{
		UnmarshalKey:   unmarshalerConfig.UnmarshalKey,
		UnmarshalValue: maps.Clone(unmarshalerConfig.UnmarshalValue),
	}
	for _, option := range options {
		option(&d)
	}
	return d
}

func unmarshalParam(cfg UnmarshalerConfig, key string, data []byte) (any, error) {
	unmarshller, found := cfg.UnmarshalValue[key]
	if found {
		return unmarshller(data)
	}
	var value any
	err := json.Unmarshal(data, &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func ConfigureUnmarshaler(f func(d *UnmarshalerConfig)) {
	f(&unmarshalerConfig)
}

type UnmarshalerConfig struct {
	UnmarshalKey   func(key string) string
	UnmarshalValue map[string]func(value []byte) (any, error)
}

var unmarshalerConfig = UnmarshalerConfig{
	UnmarshalKey: func(name string) string {
		r, n := utf8.DecodeRuneInString(name)
		if unicode.IsUpper(r) {
			return name
		}
		return string(unicode.ToUpper(r)) + name[n:]
	},
	UnmarshalValue: map[string]func(data []byte) (any, error){
		keyCode: func(data []byte) (any, error) {
			var code errors.Code
			err := json.Unmarshal(data, &code)
			if err != nil {
				return nil, err
			}
			return code, nil
		},
		keyMessage: func(data []byte) (any, error) {
			var message string
			err := json.Unmarshal(data, &message)
			if err != nil {
				return nil, err
			}
			return message, nil
		},
		keyCause: func(data []byte) (any, error) {
			e := &Error{}
			err := json.Unmarshal(data, e)
			if err == nil {
				return e, nil
			}
			var str string
			err = json.Unmarshal(data, &str)
			if err == nil {
				return stderrors.New(str), nil
			}
			return stderrors.New(string(data)), nil
		},
	},
}

type keyUnmarshalError struct {
	key string
	err error
}

func (e keyUnmarshalError) Error() string {
	return fmt.Sprintf("unmarshal %s: %s", e.key, e.err.Error())
}

func (e keyUnmarshalError) Unwrap() error {
	return e.err
}

func (e keyUnmarshalError) Cause() error {
	return e.err
}

type keyMissingError struct {
	key string
}

func (e keyMissingError) Error() string {
	return fmt.Sprintf("missing %s", e.key)
}

type keyCastError struct {
	key string
}

func (e keyCastError) Error() string {
	return fmt.Sprintf("cast %s unmarshaled value", e.key)
}
