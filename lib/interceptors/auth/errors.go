package auth

import (
	"errors"
)

var (
	ErrUnAuthorization = errors.New("Unauthorized")
	ErrInternal        = errors.New("Internal error")
)
