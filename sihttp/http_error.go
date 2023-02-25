package sihttp

import (
	"bytes"
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
	msg := bytes.Buffer{}
	if e.Response == nil {
		msg.WriteString("status: unknown")
	} else {
		msg.WriteString(fmt.Sprintf("status: %s", e.Response.Status))
	}

	if e.Body != nil {
		if msg.Len() > 0 {
			msg.WriteString(", ")
		}
		msg.WriteString(fmt.Sprintf("body: %s", e.Body))
	}
	return msg.String()
}

func (e Error) GetStatusCode() int {
	if e.Response != nil {
		return e.Response.StatusCode
	}
	return http.StatusInternalServerError
}
func (e Error) GetStatus() string {
	if e.Response != nil {
		return e.Response.Status
	}
	return http.StatusText(http.StatusInternalServerError)
}
