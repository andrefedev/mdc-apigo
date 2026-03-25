package okgrpcx

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"apigo/internal/features/auth"
	"apigo/internal/features/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type logRecord struct {
	level slog.Level
	attrs map[string]any
}

type captureHandler struct {
	records []logRecord
}

func (h *captureHandler) Enabled(context.Context, slog.Level) bool { return true }

func (h *captureHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := make(map[string]any)
	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value.Any()
		return true
	})

	h.records = append(h.records, logRecord{
		level: record.Level,
		attrs: attrs,
	})
	return nil
}

func (h *captureHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }

func (h *captureHandler) WithGroup(_ string) slog.Handler { return h }

func TestUnaryLoggingInterceptorLogsAppErrorFields(t *testing.T) {
	handler := &captureHandler{}
	prev := slog.Default()
	slog.SetDefault(slog.New(handler))
	t.Cleanup(func() { slog.SetDefault(prev) })

	appErr := fmt.Errorf("Auth.Code: %w", auth.ErrInvalidPhone)

	_, err := UnaryLoggingInterceptor(
		context.Background(),
		struct{}{},
		&grpc.UnaryServerInfo{FullMethod: "/muydelcampo.v1.AuthService/Code"},
		func(ctx context.Context, req any) (any, error) {
			return nil, appErr
		},
	)
	if err == nil {
		t.Fatal("expected error")
	}

	if len(handler.records) != 1 {
		t.Fatalf("expected 1 log record, got %d", len(handler.records))
	}

	record := handler.records[0]
	if record.level != slog.LevelWarn {
		t.Fatalf("expected warn level, got %v", record.level)
	}
	if got := record.attrs["grpc_code"]; got != "Unknown" {
		t.Fatalf("expected grpc_code Unknown, got %#v", got)
	}
	if got := record.attrs["app_kind"]; got != kindValidation {
		t.Fatalf("expected app_kind %q, got %#v", kindValidation, got)
	}
	if got := record.attrs["app_code"]; got != "auth.invalid_phone" {
		t.Fatalf("expected app_code auth.invalid_phone, got %#v", got)
	}
}

func TestUnaryLoggingInterceptorExtractsAppFieldsFromStatusDetails(t *testing.T) {
	handler := &captureHandler{}
	prev := slog.Default()
	slog.SetDefault(slog.New(handler))
	t.Cleanup(func() { slog.SetDefault(prev) })

	grpcErr := StatusError(fmt.Errorf("Users.Me: %w", users.ErrAuthenticationRequired))

	_, err := UnaryLoggingInterceptor(
		context.Background(),
		struct{}{},
		&grpc.UnaryServerInfo{FullMethod: "/muydelcampo.v1.UserService/UserMe"},
		func(ctx context.Context, req any) (any, error) {
			return nil, grpcErr
		},
	)
	if err == nil {
		t.Fatal("expected error")
	}

	if len(handler.records) != 1 {
		t.Fatalf("expected 1 log record, got %d", len(handler.records))
	}

	record := handler.records[0]
	if got := record.attrs["grpc_code"]; got != codes.Unauthenticated.String() {
		t.Fatalf("expected grpc_code %q, got %#v", codes.Unauthenticated.String(), got)
	}
	if got := record.attrs["app_kind"]; got != kindUnauthorized {
		t.Fatalf("expected app_kind %q, got %#v", kindUnauthorized, got)
	}
	if got := record.attrs["app_code"]; got != "users.authentication_required" {
		t.Fatalf("expected app_code users.authentication_required, got %#v", got)
	}
}
