package inventory

import "errors"

var (
	ErrServerNotFound = errors.New("server not found")
	ErrDuplicateServer = errors.New("server with details already exists")
)
