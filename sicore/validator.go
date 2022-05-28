package sicore

import "io"

// ReadValidator validates read bytes
type ReadValidator interface {
	validate([]byte, error) (bool, error)
}

// validateFunc wraps a function that conforms to ReadValidator interface
type ValidateFunc func([]byte, error) (bool, error)

// implements ReadValidator's validate method
func (v ValidateFunc) validate(b []byte, errIn error) (bool, error) {
	return v(b, errIn)
}

// DefaultValidator simply checks EOF
func DefaultValidator() ReadValidator {
	return ValidateFunc(func(b []byte, errIn error) (bool, error) {
		if errIn != nil {
			if errIn == io.EOF {
				return true, nil
			}
			return false, errIn
		}

		return false, nil
	})
}

// func defaultValidate(b []byte, errIn error) (bool, error) {
// 	if errIn != nil {
// 		if errIn == io.EOF {
// 			return true, nil
// 		}
// 		return false, errIn
// 	}

// 	return false, nil
// }
