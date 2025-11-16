package http

import (
	"encoding/json"

	goerrors "errors"

	"pr-reviewer-assign-service/pkg/errors"
	jsonerr "pr-reviewer-assign-service/pkg/errors/json"
)

const keyStatus = "HttpStatus"

func WithStatus(status int) errors.Param {
	return errors.Param{Name: keyStatus, Value: status}
}

func GetStatus(err error) (int, bool) {
	var e *errors.Error
	ok := goerrors.As(err, &e)
	if !ok {
		return 0, false
	}

	status, ok := e.Get(keyStatus).(int)
	if !ok {
		return 0, false
	}

	return status, ok
}

func init() {
	jsonerr.ConfigureMarshaler(func(e *jsonerr.MarshalerConfig) {
		e.MarshalValue[keyStatus] = func(value any) ([]byte, error) {
			return json.Marshal(value)
		}
	})
	jsonerr.ConfigureUnmarshaler(func(d *jsonerr.UnmarshalerConfig) {
		d.UnmarshalValue[keyStatus] = func(data []byte) (any, error) {
			var status int

			err := json.Unmarshal(data, &status)
			if err != nil {
				return nil, err
			}

			return status, nil
		}
	})
}

func JSONPrivate(e *jsonerr.MarshalerConfig) {
	private := e.IsPrivateKey
	e.IsPrivateKey = func(name string) bool {
		if name == keyStatus {
			return true
		}

		return private(name)
	}
}
