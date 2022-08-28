package justrpc

import "fmt"

const (
	InvalidRequestFormat = "justrpc.invalidRequestFormat"
	VersionMismatch      = "justrpc.versionMismatch"
	MethodNotFound       = "justrpc.methodNotFound"
	InvalidArgs          = "justrpc.invalidArgs"
	InternalError        = "justrpc.internalError"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func Errorf(code string, format string, a ...any) error {
	return &Error{Code: code, Message: fmt.Sprintf(format, a...)}
}
