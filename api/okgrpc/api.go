package okgrpc

import (
	"apigo/internal/app"
	"apigo/internal/modules/whatsapp/messages"
	v1 "apigo/protobuf/gen/v1"
)

type Server struct {
	//vendor *external.Vendor
	// repository *repository.Repository // REPOSITORY
	v1.UnsafeApiServiceServer

	repository     *app.Repository
	useservice     *app.UseService
	messageservice *messages.Service
	// ServerDeps
}

type ServerDeps struct {
	Repository     *app.Repository
	UseService     *app.UseService
	MessageService *messages.Service
}

func NewServer(deps ServerDeps) *Server {
	return &Server{
		repository:     deps.Repository,
		useservice:     deps.UseService,
		messageservice: deps.MessageService,
	}
}
