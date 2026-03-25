package messages

import (
	"apigo/internal/modules/whatsapp"
	"context"
	"fmt"
)

// Service sends WhatsApp Cloud API messages.
type Service struct {
	client *whatsapp.Client
}

// NewService creates a messages service.
func NewService(client *whatsapp.Client) *Service {
	return &Service{client: client}
}

func (s *Service) SendTemplate(ctx context.Context, req *TemplateMessageRequest) error {
	op := "Messages.Service.SendTemplate"

	body := newTemplateMessage(req)
	path := s.client.PhoneNumberId() + "/" + "messages"
	if err := s.client.Post(ctx, path, body); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) endpoint() string {
	return s.client.PhoneNumberId() + "/" + "messages"
}
