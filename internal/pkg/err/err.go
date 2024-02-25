package err

import "errors"

var (
	ErrInvalidEnvironment = errors.New("provided environment is invalid")
)
