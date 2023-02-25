package sihttp

import (
	"fmt"
	"net/http"
)

type Error struct {
	Response *http.Response
	Body     []byte
}

// func NewError(status int, message string) *Error {
// 	return &Error{
// 		Status:  status,
// 		Message: message,
// 	}
// }

func (e Error) Error() string {
	if e.Response == nil {
		return "status: unknown"
	}
	return fmt.Sprintf("status: %s, body: %s", e.Response.Status, e.Body)
}
