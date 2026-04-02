package model

import "errors"

var (
	ErrPartNotFound    = errors.New("part not found")
	ErrInvalidUUID     = errors.New("invalid uuid")
	ErrInvalidCategory = errors.New("invalid category")
)
