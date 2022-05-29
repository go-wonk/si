package sicore

import "errors"

const maxInt = int(^uint(0) >> 1)

var ErrTooLarge = errors.New("buf too large")

// grow a byte slice's capacity or allocate more space(len) by n
func grow(b *[]byte, n int) (int, error) {
	c := cap(*b)
	l := len(*b)
	a := c - l // available
	if n <= a {
		*b = (*b)[:l+n]
		return l, nil
	}

	if l+n <= c {
		// if needed length is lte c
		return l, nil
	}

	if c > maxInt-c-n {
		// too large
		return l, ErrTooLarge
	}

	newBuf := make([]byte, c*2+n)
	copy(newBuf, (*b)[0:])
	*b = newBuf[:l+n]
	return l, nil
}

// growCap grows the capacity of byte slice by n
func growCap(b *[]byte, n int) error {
	c := cap(*b)
	l := len(*b)
	a := c - l // available
	if n <= a {
		return nil
	}

	if l+n <= c {
		// if needed length is lte c
		return nil
	}

	if c > maxInt-c-n {
		// too large
		return ErrTooLarge
	}

	newBuf := make([]byte, c+n)
	copy(newBuf, (*b)[0:])
	*b = newBuf[:l]
	return nil
}
