package api

import (
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	auth_v1 "github.com/ChopX4/raketka/shared/pkg/proto/auth/v1"
)

func validateRegisterRequest(req *auth_v1.RegisterRequest) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}

	if strings.TrimSpace(req.GetLogin()) == "" {
		return status.Error(codes.InvalidArgument, "login is required")
	}

	if strings.TrimSpace(req.GetPassword()) == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if strings.TrimSpace(req.GetEmail()) == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if _, err := mail.ParseAddress(req.GetEmail()); err != nil {
		return status.Error(codes.InvalidArgument, "email is invalid")
	}

	for _, method := range req.GetNotificationMethods() {
		if method == nil {
			return status.Error(codes.InvalidArgument, "notification method is invalid")
		}

		if strings.TrimSpace(method.GetProviderName()) == "" {
			return status.Error(codes.InvalidArgument, "notification method provider_name is required")
		}

		if strings.TrimSpace(method.GetTarget()) == "" {
			return status.Error(codes.InvalidArgument, "notification method target is required")
		}
	}

	return nil
}

func validateLoginRequest(req *auth_v1.LoginRequest) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}

	if strings.TrimSpace(req.GetLogin()) == "" {
		return status.Error(codes.InvalidArgument, "login is required")
	}

	if strings.TrimSpace(req.GetPassword()) == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func validateWhoamiRequest(req *auth_v1.WhoamiRequest) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}

	if err := validateUUID(req.GetSessionUuid(), "session_uuid"); err != nil {
		return err
	}

	return nil
}

func validateGetUserRequest(req *auth_v1.GetUserRequest) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}

	if err := validateUUID(req.GetUserUuid(), "user_uuid"); err != nil {
		return err
	}

	return nil
}

func validateUUID(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return status.Errorf(codes.InvalidArgument, "%s is required", field)
	}

	if _, err := uuid.Parse(value); err != nil {
		return status.Errorf(codes.InvalidArgument, "%s is invalid", field)
	}

	return nil
}
