package main

import (
	"apigo/api/okgrpc"
	"apigo/internal/features/users"
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
	"apigo/internal/platforms/configx"
	"apigo/internal/platforms/loggerx"
	"apigo/internal/platforms/okgrpcx"
	muydelcampov1 "apigo/protobuf/gen/v1"
)

func main() {
	cfg, err := configx.Load()
	if err != nil {
		slog.Error("grpc server main: load configx", "err", err)
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

	serverx := okgrpc.NewServer(
		okgrpc.ServerDeps{
			AuthRepository: auth.NewRepository(pgdb),
			UserRepository: users.NewRepository(pgdb),
			WhatsAppClient: whatsapp.NewClient(
				whatsapp.Config{
					ApiToken: cfg.WhatsAppToken,
					ApiPhone: cfg.WhatsAppPhone,
				},
			),
		},
	)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			okgrpcx.UnaryErrorInterceptor,
			okgrpcx.UnaryLoggingInterceptor,
			okgrpc.SessionUnaryInterceptor(serverx),
			okgrpc.AuthorizeUnaryInterceptor(serverx),
		),
	)

	muydelcampov1.RegisterAuthServiceServer(
		grpcServer,
		okgrpc.NewAuthService(serverx),
	)

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
		if err := grpcServer.Serve(lis); err != nil {
			serverErr <- err
		}
	}()

	select {
	case <-signalCtx.Done():
		slog.Info("grpc shutdown requested")
		grpcServer.GracefulStop()
	case err := <-serverErr:
		slog.Error("grpc server serve", "err", err)
		os.Exit(1)
	}

	slog.Info("grpc server stopped")
}
