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
	DioTimeout   time.Duration
	DioUserAgent string
}

// WithDefaults fills empty optional values.
func (c Config) WithDefaults() Config {
	if c.ApiVersion == "" {
		c.ApiVersion = DefaultVersion
	}
	if c.ApiBaseUrl == "" {
		c.ApiBaseUrl = DefaultBaseUrl
	}
	if c.DioTimeout == 0 {
		c.DioTimeout = 15 * time.Second
	}
	if c.DioUserAgent == "" {
		c.DioUserAgent = "apigo/run"
	}
	return c
}
