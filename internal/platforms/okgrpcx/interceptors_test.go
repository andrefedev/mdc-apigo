package okgrpcx

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"apigo/internal/apperr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func TestUnaryLoggingInterceptorLogsSemanticError(t *testing.T) {
	handler := &captureHandler{}
	prev := slog.Default()
	slog.SetDefault(slog.New(handler))
	t.Cleanup(func() { slog.SetDefault(prev) })

	appErr := fmt.Errorf("Auth.Code: %w", apperr.ErrInvalidPhone)

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
	if got := record.attrs["grpc_code"]; got != codes.InvalidArgument.String() {
		t.Fatalf("expected grpc_code %q, got %#v", codes.InvalidArgument.String(), got)
	}
	if _, ok := record.attrs["err"]; !ok {
		t.Fatal("expected err attribute")
	}
	if got := fmt.Sprint(record.attrs["err"]); got != "Auth.Code: invalid phone" {
		t.Fatalf("expected wrapped error string, got %q", got)
	}
}

func TestStatusErrorMapsSemanticError(t *testing.T) {
	err := StatusError(fmt.Errorf("Auth.Code: %w", apperr.ErrInvalidPhone))
	st := status.Convert(err)

	if st.Code() != codes.InvalidArgument {
		t.Fatalf("expected code %q, got %q", codes.InvalidArgument, st.Code())
	}
	if st.Message() != "El número de teléfono no es válido" {
		t.Fatalf("expected public message, got %q", st.Message())
	}
}
