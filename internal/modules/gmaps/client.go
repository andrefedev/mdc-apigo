package gmaps

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
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

func ParseSessionToken(token string) (maps.PlaceAutocompleteSessionToken, error) {
	const op = "Mapx.Client.ParseSessionToken"

	t, err := uuid.Parse(token)
	if err != nil {
		return maps.PlaceAutocompleteSessionToken{}, fmt.Errorf("%s: %w", op, ErrPlaceTokenInvalid)
	}

	return maps.PlaceAutocompleteSessionToken(t), nil
}
