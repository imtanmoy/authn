package errorx

import "errors"

var (
	ErrorNotFound     = errors.New("resource not found")
	ErrTokenExpired   = errors.New("token expired")
	ErrInternalDB     = errors.New("internal database error")
	ErrInternalServer = errors.New("internal server error")
)
