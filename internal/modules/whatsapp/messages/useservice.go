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

func (s *Service) SendTemplate(ctx context.Context, req *SendTemplateMessage) error {
	op := "Messages.Service.SendTemplate"

	body := NewTemplateMessageReq(req)
	if err := s.client.Post(ctx, s.endpoint(), body); err != nil {
		return apperr.Internal(op, err)
	}

	return nil

	// if err := s.client.Post(ctx, s.phoneNumberID+"/messages", request, response); err != nil {
	// 	return nil, err
	// }
	// return response, nil
}

func (s *Service) endpoint() string {
	return s.client.PhoneNumberId() + "/" + "messages"
}
