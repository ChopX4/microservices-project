package converter

import (
	"github.com/ChopX4/raketka/iam/internal/model"
	repoModel "github.com/ChopX4/raketka/iam/internal/repository/model"
)

func UserToRepo(user model.User) repoModel.User {
	return repoModel.User{
		Uuid:                user.Uuid,
		Login:               user.Login,
		Email:               user.Email,
		HashPassword:        user.HashPassword,
		NotificationMethods: NotificationMethodsToRepo(user.NotificationMethods),
	}
}

func UserToModel(user repoModel.User) model.User {
	return model.User{
		Uuid:                user.Uuid,
		Login:               user.Login,
		Email:               user.Email,
		HashPassword:        user.HashPassword,
		NotificationMethods: NotificationMethodsToModel(user.NotificationMethods),
	}
}

func NotificationMethodToRepo(method model.NotificationMethod) repoModel.NotificationMethod {
	return repoModel.NotificationMethod{
		ProviderName: method.ProviderName,
		Target:       method.Target,
	}
}

func NotificationMethodToModel(method repoModel.NotificationMethod) model.NotificationMethod {
	return model.NotificationMethod{
		ProviderName: method.ProviderName,
		Target:       method.Target,
	}
}

func NotificationMethodsToRepo(methods []model.NotificationMethod) []repoModel.NotificationMethod {
	if len(methods) == 0 {
		return nil
	}

	result := make([]repoModel.NotificationMethod, 0, len(methods))
	for _, method := range methods {
		result = append(result, NotificationMethodToRepo(method))
	}

	return result
}

func NotificationMethodsToModel(methods []repoModel.NotificationMethod) []model.NotificationMethod {
	if len(methods) == 0 {
		return nil
	}

	result := make([]model.NotificationMethod, 0, len(methods))
	for _, method := range methods {
		result = append(result, NotificationMethodToModel(method))
	}

	return result
}
