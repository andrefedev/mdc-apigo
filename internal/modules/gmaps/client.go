package gmaps

import (
	"fmt"
	"strings"

	"googlemaps.github.io/maps"
)

type Client struct {
	config Config
	client *maps.Client
}

func NewClient(apiKey string) (*Client, error) {
	const op = "gmaps.NewClient"

	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, fmt.Errorf("%s: %w", op, ErrApiKeyRequired)
	}

	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		client: client,
		config: DefaultConfig(),
	}, nil
}

func NewSessionToken() maps.PlaceAutocompleteSessionToken {
	return maps.NewPlaceAutocompleteSessionToken()
}
