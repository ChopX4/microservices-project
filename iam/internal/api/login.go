package api

import (
	"context"

	"github.com/ChopX4/raketka/iam/internal/converter"
	auth_v1 "github.com/ChopX4/raketka/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, req *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {
	if err := validateLoginRequest(req); err != nil {
		return nil, err
	}

	sessionUUID, err := a.iamService.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, mapError(err)
	}

	return converter.LoginResponseToProto(sessionUUID), nil
}
