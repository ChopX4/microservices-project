package service

import (
	"context"

	"github.com/ChopX4/raketka/iam/internal/model"
)

type IamService interface {
	Login(ctx context.Context, login, password string) (string, error)
	Register(ctx context.Context, req model.RegisterRequest) (string, error)
	Whoami(ctx context.Context, sessionUUID string) (model.WhoamiResponse, error)
	GetUser(ctx context.Context, userUUID string) (model.GetUserResponse, error)
}
