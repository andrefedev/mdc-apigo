package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"apigo/api/okgrpc"
	"apigo/internal/app"
	"apigo/internal/modules/gmaps"
	"apigo/internal/modules/postgres"
	"apigo/internal/modules/whatsapp"
	"apigo/internal/modules/whatsapp/messages"
	"apigo/internal/platforms/confx"
	"apigo/internal/platforms/loggerx"
	v1 "apigo/protobuf/gen/v1"
)

func main() {
	cfg, err := confx.Load()
	if err != nil {
		slog.Error("grpc server main: load confx", "err", err)
		os.Exit(1)
	}

	loggerx.SetupLogger(cfg.Env)

	// ############
	// # DATABASE #
	// ############
	ctx := context.Background()
	pool, err := postgres.Open(ctx, cfg.PgDatabaseUrl)
	if err != nil {
		slog.Error("grpc server main: open db", "err", err)
		os.Exit(1)
	}
	defer pool.Close()
	pgdb := postgres.NewPgdb(pool)
	// ################
	// # END_DATABASE #
	// ################

	// ############
	// # WHATSAPP #
	// ############
	waba := whatsapp.NewClient(
		whatsapp.Config{
			ApiToken: cfg.WhatsAppToken,
			ApiPhone: cfg.WhatsAppPhone,
		},
	)
	messageservice := messages.NewService(waba)
	// ################
	// # END WHATSAPP #
	// ################

	// ###############
	// # GOOGLE_MAPS #
	// ###############
	mapx, err := gmaps.NewClient(cfg.GoogleMapsApiKey)
	if err != nil {
		slog.Error("grpc server main: maps client", "err", err)
		os.Exit(1)
	}
	// ###################
	// # END_GOOGLE_MAPS #
	// ###################

	// ################
	// # APP SERVICES #
	// ################
	repo := app.NewRepository(pgdb)
	useservice := app.NewUseService(app.UseServiceDeps{
		Repository:     repo,
		GoogleMapx:     mapx,
		MessageService: messageservice,
	})
	// ####################
	// # END APP SERVICES #
	// ####################

	serverx := okgrpc.NewServer(
		okgrpc.ServerDeps{
			Repository:     repo,
			UseService:     useservice,
			MessageService: messageservice,
		},
	)

	serverg := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			okgrpc.UnaryErrorInterceptor,
			okgrpc.UnaryLoggingInterceptor,
			okgrpc.SessionUnaryInterceptor(serverx),
			okgrpc.AuthorizeUnaryInterceptor(serverx),
		),
	)

	v1.RegisterApiServiceServer(serverg, serverx)

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		slog.Error("grpc server main: listen", "addr", cfg.Port, "err", err)
		os.Exit(1)
	}

	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("grpc server listening", "addr", cfg.Port)
		if err := serverg.Serve(lis); err != nil {
			serverErr <- err
		}
	}()

	select {
	case <-signalCtx.Done():
		slog.Info("grpc shutdown requested")
		serverg.GracefulStop()
	case err := <-serverErr:
		slog.Error("grpc server serve", "err", err)
		os.Exit(1)
	}

	slog.Info("grpc server stopped")
}
