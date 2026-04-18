package api

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ChopX4/raketka/iam/internal/model"
)

func mapError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, model.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, model.ErrSessionNotFound):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, model.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, model.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return err
	}
}
