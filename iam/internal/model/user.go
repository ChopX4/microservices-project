package model

type User struct {
	Uuid                string
	Login               string
	Email               string
	HashPassword        string
	NotificationMethods []NotificationMethod
}

type RegisterRequest struct {
	Login               string
	Email               string
	Password            string
	NotificationMethods []NotificationMethod
}

type GetUserResponse struct {
	UserUUID            string
	Login               string
	Email               string
	NotificationMethods []NotificationMethod
}

type WhoamiResponse struct {
	Uuid  string
	Login string
	Email string
}

type NotificationMethod struct {
	ProviderName string
	Target       string
}
