package converter

import (
	"github.com/ChopX4/raketka/iam/internal/model"
	auth_v1 "github.com/ChopX4/raketka/shared/pkg/proto/auth/v1"
)

func RegisterRequestToModel(req *auth_v1.RegisterRequest) model.RegisterRequest {
	if req == nil {
		return model.RegisterRequest{}
	}

	return model.RegisterRequest{
		Login:               req.GetLogin(),
		Email:               req.GetEmail(),
		Password:            req.GetPassword(),
		NotificationMethods: NotificationMethodsToModel(req.GetNotificationMethods()),
	}
}

func NotificationMethodToModel(method *auth_v1.NotificationMethod) model.NotificationMethod {
	if method == nil {
		return model.NotificationMethod{}
	}

	return model.NotificationMethod{
		ProviderName: method.GetProviderName(),
		Target:       method.GetTarget(),
	}
}

func NotificationMethodsToModel(methods []*auth_v1.NotificationMethod) []model.NotificationMethod {
	if len(methods) == 0 {
		return nil
	}

	result := make([]model.NotificationMethod, 0, len(methods))
	for _, method := range methods {
		result = append(result, NotificationMethodToModel(method))
	}

	return result
}

func NotificationMethodToProto(method model.NotificationMethod) *auth_v1.NotificationMethod {
	return &auth_v1.NotificationMethod{
		ProviderName: method.ProviderName,
		Target:       method.Target,
	}
}

func NotificationMethodsToProto(methods []model.NotificationMethod) []*auth_v1.NotificationMethod {
	if len(methods) == 0 {
		return nil
	}

	result := make([]*auth_v1.NotificationMethod, 0, len(methods))
	for _, method := range methods {
		result = append(result, NotificationMethodToProto(method))
	}

	return result
}

func UserToProto(user model.User) *auth_v1.User {
	return &auth_v1.User{
		UserUuid:            user.Uuid,
		Login:               user.Login,
		Email:               user.Email,
		NotificationMethods: NotificationMethodsToProto(user.NotificationMethods),
	}
}

func WhoamiResponseToProto(response model.WhoamiResponse) *auth_v1.WhoamiResponse {
	return &auth_v1.WhoamiResponse{
		User: &auth_v1.User{
			UserUuid: response.Uuid,
			Login:    response.Login,
			Email:    response.Email,
		},
	}
}

func GetUserResponseToProto(response model.GetUserResponse) *auth_v1.GetUserResponse {
	return &auth_v1.GetUserResponse{
		User: &auth_v1.User{
			UserUuid:            response.UserUUID,
			Login:               response.Login,
			Email:               response.Email,
			NotificationMethods: NotificationMethodsToProto(response.NotificationMethods),
		},
	}
}

func LoginResponseToProto(sessionUUID string) *auth_v1.LoginResponse {
	return &auth_v1.LoginResponse{
		SessionUuid: sessionUUID,
	}
}

func RegisterResponseToProto(userUUID string) *auth_v1.RegisterResponse {
	return &auth_v1.RegisterResponse{
		UserUuid: userUUID,
	}
}
