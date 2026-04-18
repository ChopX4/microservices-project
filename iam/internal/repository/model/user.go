package model

type User struct {
	Uuid                string
	Login               string
	Email               string
	HashPassword        string
	NotificationMethods []NotificationMethod
}

type NotificationMethod struct {
	// Имя провайдера (например, telegram, email, push)
	ProviderName string
	// Адрес получателя — email, ID чата и т.д.
	Target string
}
