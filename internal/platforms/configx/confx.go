package configx

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	Env                string
	Port               string
	WhatsAppToken      string
	WhatsAppPhone      string
	PgDatabaseUrl      string
	GoogleMapsApiKey   string
	GoogleGeminiApiKey string
}

func Load() (Config, error) {
	cfg := Config{
		Env:              getenv("ENV", "dev"),
		Port:             listenAddr(getenv("PORT", "8080")),
		WhatsAppToken:    os.Getenv("WHATSAPP_TOKEN"),
		WhatsAppPhone:    os.Getenv("WHATSAPP_PHONE"),
		PgDatabaseUrl:    os.Getenv("PG_DATABASE_URL"),
		GoogleMapsApiKey: os.Getenv("GOOGLE_MAPS_API_KEY"),
	}

	if cfg.WhatsAppToken == "" {
		return Config{}, errors.New("configx load(): you must set your 'WHATSAPP_TOKEN' env var")
	}

	if cfg.WhatsAppPhone == "" {
		return Config{}, errors.New("configx load(): you must set your 'WHATSAPP_PHONE' env var")
	}

	if cfg.PgDatabaseUrl == "" {
		return Config{}, errors.New("configx load(): you must set your 'PG_DATABASE_URL' env var")
	}

	return cfg, nil
}

func getenv(k, def string) string {
	v := strings.TrimSpace(os.Getenv(k))
	if v != "" {
		return v
	}
	return def
}

func listenAddr(value string) string {
	v := strings.TrimSpace(value)
	if v == "" {
		return ":8080"
	}

	if strings.HasPrefix(v, ":") || strings.Contains(v, ":") {
		return v
	}

	return ":" + v
}
