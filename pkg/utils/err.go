package utils

import (
	"errors"
	"fmt"
)

func NewErrf(format string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(format, a...))
}
