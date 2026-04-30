package iam

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/ChopX4/raketka/iam/internal/model"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

func (s *service) Login(ctx context.Context, login, password string) (string, error) {
	if err := s.validateLoginRequest(login, password); err != nil {
		return "", err
	}

	user, err := s.userRepository.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return "", model.ErrInvalidCredentials
		}

		logger.Error(ctx, "failed to get user during login", zap.String("login", login), zap.Error(err))
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(password)); err != nil {
		return "", model.ErrInvalidCredentials
	}

	sessionUUID := uuid.New()

	if err := s.sessionRepository.Create(ctx, model.Session{SessionUUID: sessionUUID.String(), UserUUID: user.Uuid}, s.sessionTTL); err != nil {
		logger.Error(ctx, "failed to create session during login", zap.String("login", login), zap.String("user_uuid", user.Uuid), zap.Error(err))
		return "", err
	}

	return sessionUUID.String(), nil
}
