package loggerx

import (
	"log/slog"
	"os"
)

// SetupLogger configura el loggerx global de la aplicación (slog).
// Si env == "dev", usa texto plano para facilitar la lectura local.
// De lo contrario, usa JSON estructurado optimizado para Google Cloud Run.
func SetupLogger(env string) {
	var h slog.Handler

	if env == "dev" {
		// Formato de texto para desarrollo local
		h = slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		)
	} else {
		// Formato JSON para Google Cloud Run / Producción
		h = slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					// Google Cloud Logging usa "severity" en lugar de "level"
					if a.Key == slog.LevelKey {
						a.Key = "severity"
						// Mapear el nivel a las severidades de GCP (opcional, GCP suele entender INFO/ERROR de slog)
						// if level, ok := a.Value.Any().(slog.Level); ok {
						// 	if level == slog.LevelWarn {
						// 		a.Value = slog.StringValue("WARNING")
						// 	}
						// }
					}
					// Google Cloud Logging usa "message" en lugar de "msg"
					if a.Key == slog.MessageKey {
						a.Key = "message"
					}
					return a
				},
			},
		)
	}

	// Establecer el loggerx configurado como el predeterminado global
	l := slog.New(h)
	slog.SetDefault(l)
}
