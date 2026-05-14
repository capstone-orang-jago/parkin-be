package domain

type WaService interface {
	SendMessage(telephone string, message string) error
}
