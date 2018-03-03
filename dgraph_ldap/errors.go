package dgraph_ldap

import "fmt"

type ErrorCode uint

const (
	ErrorGraphNotInitialized ErrorCode = iota
)

var errorDescriptions = map[ErrorCode]string{
	ErrorGraphNotInitialized: "ErrorGraphNotInitialized",
}

type Error struct {
	code ErrorCode
	info string
}

func (e *Error) Error() string {
	d := errorDescriptions[e.code]
	if d != "" {
		d += " "
	}
	return fmt.Sprintf("[%s%d] %s", d, e.code, e.info)
}

func (e *Error) Code() ErrorCode {
	return e.code
}

func NewError(code ErrorCode, describe string) *Error {
	return &Error{code, describe}
}

