package okgrpc

import (
	"apigo/internal/features/auth"
	"apigo/internal/features/users"
	"apigo/internal/modules/whatsapp"
	"apigo/internal/modules/whatsapp/messages"
	"apigo/internal/platforms/configx"
	v1 "apigo/protobuf/gen/v1"
)

type Server struct {
	//vendor *external.Vendor
	// repository *repository.Repository // REPOSITORY
	v1.UnsafeApiServiceServer

	AuthService *auth.Service
	UserService *users.Service
	// ServerDeps
}

type ServerDeps struct {
	Config         configx.Config
	WhatsAppClient *whatsapp.Client
	AuthRepository *auth.Repository
	UserRepository *users.Repository
}

func NewServer(deps ServerDeps) *Server {
	return &Server{
		AuthService: auth.NewService(
			auth.ServiceDeps{
				Repository:     deps.AuthRepository,
				MessageService: messages.NewService(deps.WhatsAppClient),
			},
		),
		UserService: users.NewService(
			users.ServiceDeps{
				UserRepository: deps.UserRepository,
			},
		),
	}
}
