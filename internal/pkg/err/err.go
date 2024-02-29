package err

import "errors"

var (
	ErrInvalidEnvironment = errors.New("provided environment is invalid")
	ErrNoSuchKeyInCache   = errors.New("there no such key in cache")
	ErrNotFound           = errors.New("item was not found")
	ErrEmptyItems         = errors.New("items list should not be empty")
	ErrInvalidField       = errors.New("non of the fields should be empty")
	ErrAlreadyExists      = errors.New("such key already exists")
	ErrInternal           = errors.New("internal server error happend")
)
