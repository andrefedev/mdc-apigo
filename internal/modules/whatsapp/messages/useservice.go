package messages

import (
	"apigo/internal/modules/whatsapp"
	"apigo/internal/platforms/apperr"
	"context"
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
		return apperr.Internal(op, err)
	}

	return nil
}

func (s *Service) endpoint() string {
	return s.client.PhoneNumberId() + "/" + "messages"
}
