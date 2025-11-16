package json

import (
	"encoding/json"
	"fmt"
	"maps"
	"unicode"
	"unicode/utf8"

	"pr-reviewer-assign-service/pkg/errors"
)

type MarshalerOption func(e *MarshalerConfig)

func Marshal(err error, options ...MarshalerOption) (data []byte) {
	data, eerr := json.Marshal(Error{
		Error:     err,
		Marshaler: options,
	})
	if eerr != nil {
		panic(eerr)
	}

	return data
}

var _ json.Marshaler = (*Error)(nil)

func (e Error) MarshalJSON() (bytes []byte, err error) {
	if e.Error == nil {
		return marshalNilError(), nil
	}

	concreteErr, ok := (e.Error).(*errors.Error)
	if !ok {
		return marshalAbstractError(e.Error), nil
	}

	bytes, err = marshalConcreteError(concreteErr, e.Marshaler...)
	if err == nil {
		return bytes, nil
	}

	return marshalAbstractError(e.Error), nil
}

func marshalAbstractError(e error) (bytes []byte) {
	bytes, err := json.Marshal(e.Error())
	if err != nil {
		panic(err)
	}

	return bytes
}

func marshalConcreteError(e *errors.Error, options ...MarshalerOption) (bytes []byte, err error) {
	cfg := applyMarshalerOptions(options)
	data := make(map[string]json.RawMessage, 4+len(e.Params()))

	data[cfg.MarshalKey(keyCode)], err = marshalParam(cfg, keyCode, e.Code())
	if err != nil {
		return nil, keyMarshalError{keyCode, err}
	}

	data[cfg.MarshalKey(keyMessage)], err = marshalParam(cfg, keyMessage, e.Error())
	if err != nil {
		return nil, keyMarshalError{keyMessage, err}
	}

	if !cfg.IsPrivateKey(keyCause) && e.Cause() != nil {
		data[cfg.MarshalKey(keyCause)], err = marshalParam(cfg, keyCause, e.Cause())
		if err != nil {
			return nil, keyMarshalError{keyCause, err}
		}
	}

	if !cfg.IsPrivateKey(keyStackTrace) && e.StackTrace() != nil {
		data[cfg.MarshalKey(keyStackTrace)], err = marshalParam(cfg, keyStackTrace, e.StackTrace())
		if err != nil {
			return nil, keyMarshalError{keyStackTrace, err}
		}
	}

	for _, param := range e.Params() {
		if cfg.IsPrivateKey(param.Name) {
			continue
		}

		data[cfg.MarshalKey(param.Name)], err = marshalParam(cfg, param.Name, param.Value)
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(data)
}

func marshalNilError() (bytes []byte) {
	return []byte("null")
}

func applyMarshalerOptions(options []MarshalerOption) MarshalerConfig {
	cfg := MarshalerConfig{
		IsPrivateKey: marshalerConfig.IsPrivateKey,
		MarshalKey:   marshalerConfig.MarshalKey,
		MarshalValue: maps.Clone(marshalerConfig.MarshalValue),
	}
	for _, option := range options {
		option(&cfg)
	}

	return cfg
}

func marshalParam(cfg MarshalerConfig, key string, value any) ([]byte, error) {
	marshaller, found := cfg.MarshalValue[key]
	if found {
		return marshaller(value)
	}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return jsonValue, nil
}

func ConfigureMarshaler(f func(e *MarshalerConfig)) {
	f(&marshalerConfig)
}

type MarshalerConfig struct {
	IsPrivateKey func(key string) bool
	MarshalKey   func(key string) string
	MarshalValue map[string]func(value any) ([]byte, error)
}

var marshalerConfig = MarshalerConfig{
	IsPrivateKey: func(name string) bool {
		r, _ := utf8.DecodeRuneInString(name)

		return unicode.IsLower(r)
	},
	MarshalKey: func(name string) string {
		r, n := utf8.DecodeRuneInString(name)
		if unicode.IsLower(r) {
			return name
		}

		return string(unicode.ToLower(r)) + name[n:]
	},
	MarshalValue: map[string]func(v any) ([]byte, error){
		keyCode: func(v any) ([]byte, error) {
			return json.Marshal(v)
		},
		keyMessage: func(v any) ([]byte, error) {
			return json.Marshal(v)
		},
		keyCause: func(v any) ([]byte, error) {
			e, ok := v.(*Error)
			if ok {
				return json.Marshal(e)
			}

			err, ok := v.(error)
			if ok {
				return json.Marshal(err.Error())
			}

			return json.Marshal(v)
		},
		keyStackTrace: func(v any) ([]byte, error) {
			st, ok := v.(errors.StackTrace)
			if !ok {
				return json.Marshal(v)
			}

			strs := make([]string, 0, len(st))
			for _, frame := range st {
				strs = append(
					strs,
					fmt.Sprintf("%s %s:%d", frame.Func(), frame.File(), frame.Line()),
				)
			}

			return json.Marshal(strs)
		},
	},
}

type keyMarshalError struct {
	key string
	err error
}

func (e keyMarshalError) Error() string {
	return fmt.Sprintf("marshal %s: %s", e.key, e.err.Error())
}

func (e keyMarshalError) Unwrap() error {
	return e.err
}

func (e keyMarshalError) Cause() error {
	return e.err
}
