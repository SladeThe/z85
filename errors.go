package z85

import (
	"fmt"
)

type InsufficientDestinationLengthError struct {
	want, got int
}

func (e InsufficientDestinationLengthError) Error() string {
	return fmt.Sprintf("z85: insufficient destination length: %d < %d", e.got, e.want)
}

type InvalidEncodedLengthError int

func (e InvalidEncodedLengthError) Error() string {
	return fmt.Sprintf("z85: invalid encoded length: %d", int(e))
}

type InvalidEncodedByteError byte

func (e InvalidEncodedByteError) Error() string {
	return fmt.Sprintf("z85: invalid encoded byte: %#U", rune(e))
}

type InvalidPostfixError byte

func (e InvalidPostfixError) Error() string {
	if e == 0 {
		return "z85: invalid postfix"
	}

	return fmt.Sprintf("z85: invalid postfix: %#U", rune(e))
}
