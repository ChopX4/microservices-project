package iam

import (
	"context"

	"github.com/ChopX4/raketka/iam/internal/model"
)

func (s *service) GetUser(ctx context.Context, userUUID string) (model.GetUserResponse, error) {
	user, err := s.userRepository.Get(ctx, userUUID)
	if err != nil {
		return model.GetUserResponse{}, err
	}

	return model.GetUserResponse{
		UserUUID:            user.Uuid,
		Login:               user.Login,
		Email:               user.Email,
		NotificationMethods: user.NotificationMethods,
	}, nil
}
