package errorx

import "errors"

var (
	ErrorNotFound   = errors.New("resource not found")
	ErrTokenExpired = errors.New("token expired")
)
