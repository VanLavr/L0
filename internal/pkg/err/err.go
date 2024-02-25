package err

import "errors"

var (
	ErrInvalidEnvironment = errors.New("provided environment is invalid")
	ErrNoSuchKeyInCache   = errors.New("there no such key in cache")
	ErrNotFound           = errors.New("item was not found")
)
