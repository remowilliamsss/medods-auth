package app

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrAuthHeader = errors.New("invalid auth header")
	ErrWrongToken = errors.New("wrong refresh token")
)
