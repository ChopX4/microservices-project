package api

import (
	"context"

	"github.com/ChopX4/raketka/iam/internal/converter"
	auth_v1 "github.com/ChopX4/raketka/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, req *auth_v1.WhoamiRequest) (*auth_v1.WhoamiResponse, error) {
	if err := validateWhoamiRequest(req); err != nil {
		return nil, err
	}

	modelWhoami, err := a.iamService.Whoami(ctx, req.GetSessionUuid())
	if err != nil {
		return nil, mapError(err)
	}

	return converter.WhoamiResponseToProto(modelWhoami), nil
}
