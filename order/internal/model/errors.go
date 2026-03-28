package model

import "errors"

var (
	ErrNotFound      = errors.New("order not found")
	ErrConflict      = errors.New("conflict")
	ErrBadRequest    = errors.New("bad request - validation error")
	ErrAlreadyExists = errors.New("order already exists")
)
