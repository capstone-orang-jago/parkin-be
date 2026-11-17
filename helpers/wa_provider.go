package helpers

import (
	"context"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"go.mau.fi/whatsmeow"
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

type waService struct {
	client *whatsmeow.Client
}

func NewWaService(client *whatsmeow.Client) domain.WaService {
	return &waService{client: client}
}

func (w *waService) SendMessage(telephone string, message string) error {
	formattedNumber, err := NormalizeTelephone(telephone)
	if err != nil {
		return err
	}

	jid := types.NewJID(formattedNumber, types.DefaultUserServer)

	_, err = w.client.SendMessage(context.Background(), jid, &waE2E.Message{
		Conversation: proto.String(message),
	})

	return err
}
