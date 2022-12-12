package sihttp

import "fmt"

type SiHttpError struct {
	Status  int
	Message string
}

func NewSiHttpError(status int, message string) *SiHttpError {
	return &SiHttpError{
		Status:  status,
		Message: message,
	}
}

func (e SiHttpError) Error() string {
	return fmt.Sprintf("status: %d, message: %s", e.Status, e.Message)
}
