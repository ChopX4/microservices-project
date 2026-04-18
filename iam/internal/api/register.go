package api

import (
	"context"

	"github.com/ChopX4/raketka/iam/internal/converter"
	auth_v1 "github.com/ChopX4/raketka/shared/pkg/proto/auth/v1"
)

func (a *api) Register(ctx context.Context, req *auth_v1.RegisterRequest) (*auth_v1.RegisterResponse, error) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	userUUID, err := a.iamService.Register(ctx, converter.RegisterRequestToModel(req))
	if err != nil {
		return nil, mapError(err)
	}

	return converter.RegisterResponseToProto(userUUID), nil
}
