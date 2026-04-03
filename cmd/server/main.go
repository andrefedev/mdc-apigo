package main

import (
	"apigo/api/okgrpc"
	"apigo/internal/app"
	"apigo/internal/features/users"
	"apigo/internal/modules/whatsapp/messages"
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"apigo/internal/features/auth"
	"apigo/internal/modules/postgres"
	"apigo/internal/modules/whatsapp"
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
	waba := whatsapp.NewClient(
		whatsapp.Config{
			ApiToken: cfg.WhatsAppToken,
			ApiPhone: cfg.WhatsAppPhone,
		},
	)

	repo := app.NewRepository(pgdb)
	service := app.NewService(app.ServiceDeps{
		Repository:     repo,
		MessageService: messages.NewService(waba),
	})

	serverx := okgrpc.NewServer(
		okgrpc.ServerDeps{
			Repository: repo,
			Service:    service,
			MessageC: whatsapp.NewClient(
				whatsapp.Config{
					ApiToken: cfg.WhatsAppToken,
					ApiPhone: cfg.WhatsAppPhone,
				},
			),
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
