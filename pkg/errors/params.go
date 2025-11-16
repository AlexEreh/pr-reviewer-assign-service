package errors

type Params []Param

type Param struct {
	Name  string
	Value any
}

func (p Params) toMap() map[string]any {
	m := make(map[string]any, len(p))
	for _, param := range p {
		m[param.Name] = param.Value
	}

	return m
}

const keyValidationErrors = "Errors"

func WithValidationErrors(errs map[string]string) Param {
	return Param{Name: keyValidationErrors, Value: errs}
}

func GetValidationErrors(err error) (map[string]string, bool) {
	e, ok := err.(*Error)
	if !ok {
		return nil, false
	}

	errs, ok := e.Get(keyValidationErrors).(map[string]string)

	return errs, ok
}
