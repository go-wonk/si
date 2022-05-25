package sicore

import "io"

// ReadValidator validates read bytes
type ReadValidator interface {
	validate([]byte, error) (bool, error)
}

// validateFunc wraps a function that conforms to ReadValidator interface
type validateFunc func([]byte, error) (bool, error)

// implements ReadValidator's validate method
func (v validateFunc) validate(b []byte, errIn error) (bool, error) {
	return v(b, errIn)
}

// defaultValidate simply checks EOF
func defaultValidate(b []byte, errIn error) (bool, error) {
	if errIn != nil {
		if errIn == io.EOF {
			return true, nil
		}
		return false, errIn
	}

	return false, nil
}
