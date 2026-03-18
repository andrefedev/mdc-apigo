package whatsapp

import (
	"time"
)

// Config configures the Graph API client.
type Config struct {
	ApiToken   string
	ApiPhone   string
	ApiVersion string
	ApiBaseUrl string
	// ApiRequest   *http.Client
	// DioClient    *http.Client
	HttpTimeout time.Duration
}

// WithDefaults fills empty optional values.
func (c Config) WithDefaults() Config {
	if c.ApiVersion == "" {
		c.ApiVersion = DefaultVersion
	}
	if c.ApiBaseUrl == "" {
		c.ApiBaseUrl = DefaultBaseUrl
	}
	if c.HttpTimeout == 0 {
		c.HttpTimeout = 15 * time.Second
	}
	return c
}
