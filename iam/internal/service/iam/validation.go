package iam

import (
	"net/mail"
	"strings"

	"github.com/google/uuid"

	"github.com/ChopX4/raketka/iam/internal/model"
)

const minPasswordLength = 8

func (s *service) validateRegisterRequest(req model.RegisterRequest) error {
	if strings.TrimSpace(req.Login) == "" {
		return model.ErrBadRequest
	}

	if strings.TrimSpace(req.Password) == "" {
		return model.ErrBadRequest
	}

	if len(req.Password) < minPasswordLength {
		return model.ErrBadRequest
	}

	if strings.TrimSpace(req.Email) == "" {
		return model.ErrBadRequest
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return model.ErrBadRequest
	}

	for _, method := range req.NotificationMethods {
		if strings.TrimSpace(method.ProviderName) == "" {
			return model.ErrBadRequest
		}

		if strings.TrimSpace(method.Target) == "" {
			return model.ErrBadRequest
		}
	}

	return nil
}

func (s *service) validateLoginRequest(login, password string) error {
	if strings.TrimSpace(login) == "" {
		return model.ErrBadRequest
	}

	if strings.TrimSpace(password) == "" {
		return model.ErrBadRequest
	}

	return nil
}

func (s *service) validateWhoamiRequest(sessionUUID string) error {
	return validateUUID(sessionUUID)
}

func (s *service) validateGetUserRequest(userUUID string) error {
	return validateUUID(userUUID)
}

func validateUUID(value string) error {
	if strings.TrimSpace(value) == "" {
		return model.ErrBadRequest
	}

	if _, err := uuid.Parse(value); err != nil {
		return model.ErrBadRequest
	}

	return nil
}
