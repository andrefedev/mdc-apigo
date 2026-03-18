package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"apigo/internal/features/app"
	"apigo/internal/features/auth"
	"apigo/internal/features/users"
	"apigo/internal/platforms/configx"
	"apigo/internal/platforms/loggerx"
	"apigo/internal/platforms/okhttpx"

	"apigo/internal/modules/postgres"
	"apigo/internal/modules/whatsapp"
	"apigo/internal/modules/whatsapp/messages"
)

func main() {
	cfg, err := configx.Load()
	if err != nil {
		slog.Error("server main: load configx", "err", err)
		os.Exit(1)
	}

	loggerx.SetupLogger(cfg.Env)

	ctx := context.Background()
	pool, err := postgres.Open(ctx, cfg.PgDatabaseUrl)
	if err != nil {
		slog.Error("server main: open db", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	pgdb := postgres.NewPgdb(pool)
	authRepo := auth.NewRepository(pgdb)
	userRepo := users.NewRepository(pgdb)

	// WABA CLIENT
	wabacli := whatsapp.NewClient(
		whatsapp.Config{
			ApiToken: cfg.WhatsAppToken,
			ApiPhone: cfg.WhatsAppPhone,
		},
	)

	// Modules External
	msgService := messages.NewService(wabacli)

	authService := auth.NewService(
		auth.ServiceDeps{
			AuthRepository: authRepo,
			MessageService: msgService,
		},
	)

	userService := users.NewService(
		users.ServiceDeps{
			UserRepository: userRepo,
		},
	)

	identityMiddleware := auth.NewMiddleware(authService)
	router := okhttpx.NewAppRouter(
		okhttpx.AppRouterDeps{
			AppHandler: app.NewHandler(
				app.HandlerDeps{
					// empty
				},
			),
			AuthHandler: auth.NewHandler(
				auth.HandlerDeps{
					Service: authService,
				},
			),
			UserHandler: users.NewHandler(
				users.HandlerDeps{
					Service:  userService,
					Identity: identityMiddleware,
				},
			),
			ReadyHandler: okhttpx.Readyz(pool),
		},
	)

	srv := okhttpx.NewServer(
		okhttpx.ServerConfig{
			Addr:    cfg.Port,
			Handler: router,
		},
	)

	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("server listening", "addr", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case <-signalCtx.Done():
		slog.Info("shutdown requested")
	case err := <-serverErr:
		slog.Error("server listen and serve", "err", err)
		os.Exit(1)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown", "err", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
