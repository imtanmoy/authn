package errorx

import "errors"

var (
	ErrorNotFound     = errors.New("resource not found")
	ErrTokenExpired   = errors.New("token expired")
	ErrInternalDB     = errors.New("internal database error")
	ErrInternalServer = errors.New("internal server error")
	// ErrInvalidPassword invalid password error
	ErrInvalidPassword = errors.New("invalid password")
	// ErrInvalidAccount invalid account error
	ErrInvalidAccount = errors.New("invalid account")
	// ErrUnauthorized unauthorized error
	ErrUnauthorized = errors.New("unauthorized")
)
