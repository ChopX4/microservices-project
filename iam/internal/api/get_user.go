package api

import (
	"context"

	"github.com/ChopX4/raketka/iam/internal/converter"
	auth_v1 "github.com/ChopX4/raketka/shared/pkg/proto/auth/v1"
)

func (a *api) GetUser(ctx context.Context, req *auth_v1.GetUserRequest) (*auth_v1.GetUserResponse, error) {
	if err := validateGetUserRequest(req); err != nil {
		return nil, err
	}

	modelUser, err := a.iamService.GetUser(ctx, req.GetUserUuid())
	if err != nil {
		return nil, mapError(err)
	}

	return converter.GetUserResponseToProto(modelUser), nil
}
