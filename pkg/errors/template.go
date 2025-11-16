package errors

type Template struct {
	Code    Code
	Message string
	Params  Params
}

const CodeInternalError Code = "InternalError"

var InternalError = Template{
	Code:    CodeInternalError,
	Message: "Internal error",
	Params:  Params{},
}

const CodeNotImplemented Code = "NotImplemented"

var NotImplemented = Template{
	Code:    CodeNotImplemented,
	Message: "Not implemented",
	Params:  Params{},
}
