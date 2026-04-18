package iam

import (
	"context"

	"github.com/ChopX4/raketka/iam/internal/model"
)

func (s *service) Whoami(ctx context.Context, sessionUUID string) (model.WhoamiResponse, error) {
	userUUID, err := s.sessionRepository.GetUserUUID(ctx, sessionUUID)
	if err != nil {
		return model.WhoamiResponse{}, err
	}

	user, err := s.userRepository.Get(ctx, userUUID)
	if err != nil {
		return model.WhoamiResponse{}, err
	}

	return model.WhoamiResponse{
		Uuid:  user.Uuid,
		Login: user.Login,
		Email: user.Email,
	}, nil
}
