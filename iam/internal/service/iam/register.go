package iam

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/ChopX4/raketka/iam/internal/model"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

func (s *service) Register(ctx context.Context, req model.RegisterRequest) (string, error) {
	if err := s.validateRegisterRequest(req); err != nil {
		return "", err
	}

	userUUID := uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx, "failed to hash password", zap.Error(err))
		return "", err
	}

	user := model.User{
		Uuid:                userUUID.String(),
		Login:               req.Login,
		Email:               req.Email,
		HashPassword:        string(hashedPassword),
		NotificationMethods: req.NotificationMethods,
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return "", err
	}

	return userUUID.String(), nil
}
