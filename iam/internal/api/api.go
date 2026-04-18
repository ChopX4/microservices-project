package api

import (
	"github.com/ChopX4/raketka/iam/internal/service"
	auth_v1 "github.com/ChopX4/raketka/shared/pkg/proto/auth/v1"
)

type api struct {
	auth_v1.UnimplementedAuthServiceServer
	iamService service.IamService
}

func NewIamApi(iamService service.IamService) *api {
	return &api{
		iamService: iamService,
	}
}
