package util

import "fmt"

func OpWrap(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}
