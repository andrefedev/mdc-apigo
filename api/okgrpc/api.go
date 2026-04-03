package okgrpc

import (
	"apigo/internal/app"
	"apigo/internal/features/auth"
	"apigo/internal/features/users"
	"apigo/internal/modules/whatsapp/messages"
	"apigo/internal/platforms/confx"
	v1 "apigo/protobuf/gen/v1"
)

type Server struct {
	//vendor *external.Vendor
	// repository *repository.Repository // REPOSITORY
	v1.UnsafeApiServiceServer

	Repository     *app.Repository
	UseService     *app.UseService
	MessageService *messages.Service
	// ServerDeps
}

type ServerDeps struct {
	Ser        confx.Config
	WabaClient *waba.Client
	Repository app.Repository
}

func NewServer(deps ServerDeps) *Server {
	return &Server{
		Service: deps.AppService,
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
