package sicore

import "io"

type EofChecker interface {
	Check([]byte, error) (bool, error)
}

var DefaultEofChecker = defaultEofChecker{}

type defaultEofChecker struct{}

func (c defaultEofChecker) Check(b []byte, errIn error) (bool, error) {
	if errIn != nil {
		if errIn == io.EOF {
			return true, nil
		}
		return false, errIn
	}

	return false, nil
}
