package model

import "errors"

var ErrNotFound = errors.New("order not found")
var ErrConflict = errors.New("conflict")
var ErrBadRequest = errors.New("bad request - validation error")
