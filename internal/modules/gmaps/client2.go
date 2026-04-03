package gmaps

import (
	"apigo/internal/platforms/validatex/normalizex"
	"context"
	"errors"
	"fmt"
	"strings"

	"googlemaps.github.io/maps"
)

var (
	ErrClient2Nil          = errors.New("gmaps client is nil")
	ErrAPIKeyRequired2     = errors.New("gmaps api key required")
	ErrQueryRequired2      = errors.New("gmaps query required")
	ErrPlaceIDRequired2    = errors.New("gmaps place id required")
	ErrCoordinatesInvalid2 = errors.New("gmaps invalid coordinates")
	ErrPlaceNotFound2      = errors.New("gmaps place not found")
	ErrPlaceOutOfCoverage2 = errors.New("gmaps place outside medellin coverage")
)

type Client2 struct {
	client *maps.Client
	cfg    Client2Config
}

func NewClient(apiKey string) (*Client2, error) {
	const op = "gmaps.NewClient2"

	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, fmt.Errorf("%s: %w", op, ErrAPIKeyRequired2)
	}

	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return NewClient2WithMaps(client, DefaultClient2Config())
}

func NewClient2WithMaps(client *maps.Client, cfg Client2Config) (*Client2, error) {
	const op = "gmaps.NewClient2WithMaps"

	if client == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrClient2Nil)
	}

	cfg = normalizeClient2Config(cfg)

	return &Client2{
		client: client,
		cfg:    cfg,
	}, nil
}

func NewSessionToken2() maps.PlaceAutocompleteSessionToken {
	return maps.NewPlaceAutocompleteSessionToken()
}

func (c *Client2) Autocomplete2(ctx context.Context, input string, token maps.PlaceAutocompleteSessionToken) ([]*Prediction2, error) {
	const op = "gmaps.Client2.Autocomplete2"

	query := normalizeQuery2(input)
	if query == "" {
		return nil, fmt.Errorf("%s: %w", op, ErrQueryRequired2)
	}

	if err := c.validateClient2(op); err != nil {
		return nil, err
	}

	resp, err := c.client.PlaceAutocomplete(ctx, &maps.PlaceAutocompleteRequest{
		Input:        query,
		Location:     &c.cfg.SearchCenter,
		Origin:       &c.cfg.SearchCenter,
		Radius:       c.cfg.SearchRadiusMeters,
		Language:     c.cfg.Language,
		StrictBounds: c.cfg.StrictBounds,
		SessionToken: token,
		Components: map[maps.Component][]string{
			maps.ComponentCountry: {c.cfg.CountryCode},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	limit := c.cfg.AutocompleteLimit
	if limit <= 0 || limit > len(resp.Predictions) {
		limit = len(resp.Predictions)
	}

	results := make([]*Prediction2, 0, limit)
	for i := 0; i < limit; i++ {
		p := resp.Predictions[i]
		results = append(results, &Prediction2{
			Ref:            strings.TrimSpace(p.PlaceID),
			Desc:           strings.TrimSpace(p.Description),
			Title:          strings.TrimSpace(p.StructuredFormatting.MainText),
			Subtitle:       strings.TrimSpace(p.StructuredFormatting.SecondaryText),
			DistanceMeters: p.DistanceMeters,
			Types:          cloneStrings2(p.Types),
		})
	}

	return results, nil
}

func (c *Client2) ResolveText2(ctx context.Context, input string, token maps.PlaceAutocompleteSessionToken) (*Place2, error) {
	const op = "gmaps.Client2.ResolveText2"

	query := normalizeQuery2(input)
	if query == "" {
		return nil, fmt.Errorf("%s: %w", op, ErrQueryRequired2)
	}

	if err := c.validateClient2(op); err != nil {
		return nil, err
	}

	sawOutsideCoverage := false

	predictions, err := c.Autocomplete2(ctx, query, token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for _, prediction := range predictions {
		place, err := c.PlaceDetails2(ctx, prediction.Ref, token)
		switch {
		case err == nil:
			place.Query = query
			return place, nil
		case errors.Is(err, ErrPlaceOutOfCoverage2):
			sawOutsideCoverage = true
		case errors.Is(err, ErrPlaceNotFound2):
		default:
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	place, err := c.findPlaceFromText2(ctx, query, token)
	switch {
	case err == nil:
		place.Query = query
		return place, nil
	case errors.Is(err, ErrPlaceOutOfCoverage2):
		sawOutsideCoverage = true
	case errors.Is(err, ErrPlaceNotFound2):
	default:
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	place, err = c.geocode2(ctx, query)
	switch {
	case err == nil:
		place.Query = query
		return place, nil
	case errors.Is(err, ErrPlaceOutOfCoverage2):
		sawOutsideCoverage = true
	case errors.Is(err, ErrPlaceNotFound2):
	default:
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if sawOutsideCoverage {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceOutOfCoverage2)
	}

	return nil, fmt.Errorf("%s: %w", op, ErrPlaceNotFound2)
}

func (c *Client2) PlaceDetails2(ctx context.Context, placeID string, token maps.PlaceAutocompleteSessionToken) (*Place2, error) {
	const op = "gmaps.Client2.PlaceDetails2"

	placeID = strings.TrimSpace(placeID)
	if placeID == "" {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceIDRequired2)
	}

	if err := c.validateClient2(op); err != nil {
		return nil, err
	}

	fields := []maps.PlaceDetailsFieldMask{
		maps.PlaceDetailsFieldMaskAddressComponent,
		maps.PlaceDetailsFieldMaskFormattedAddress,
		maps.PlaceDetailsFieldMaskGeometryLocation,
		maps.PlaceDetailsFieldMaskName,
		maps.PlaceDetailsFieldMaskPlaceID,
		maps.PlaceDetailsFieldMaskTypes,
	}

	result, err := c.client.PlaceDetails(ctx, &maps.PlaceDetailsRequest{
		PlaceID:      placeID,
		Language:     c.cfg.Language,
		Region:       c.cfg.Region,
		Fields:       fields,
		SessionToken: token,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	place := c.placeFromDetails2(result)
	if !place.InCoverage {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceOutOfCoverage2)
	}

	return place, nil
}

func (c *Client2) ReverseGeocode2(ctx context.Context, lat, lng float64) (*Place2, error) {
	const op = "gmaps.Client2.ReverseGeocode2"

	if !isValidLatLng2(lat, lng) {
		return nil, fmt.Errorf("%s: %w", op, ErrCoordinatesInvalid2)
	}

	if err := c.validateClient2(op); err != nil {
		return nil, err
	}

	latLng := &maps.LatLng{Lat: lat, Lng: lng}

	results, err := c.client.ReverseGeocode(ctx, &maps.GeocodingRequest{
		LatLng:       latLng,
		Region:       c.cfg.Region,
		Language:     c.cfg.Language,
		ResultType:   []string{"street_address", "premise", "subpremise"},
		LocationType: []maps.GeocodeAccuracy{maps.GeocodeAccuracyRooftop, maps.GeocodeAccuracyRangeInterpolated},
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(results) == 0 {
		results, err = c.client.ReverseGeocode(ctx, &maps.GeocodingRequest{
			LatLng:     latLng,
			Region:     c.cfg.Region,
			Language:   c.cfg.Language,
			ResultType: []string{"route", "neighborhood", "sublocality", "locality"},
		})
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	result := c.pickBestGeocodingResult2(results)
	if result == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceNotFound2)
	}

	place := c.placeFromGeocode2(*result, "reverse_geocode")
	if !place.InCoverage {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceOutOfCoverage2)
	}

	return place, nil
}

func (c *Client2) findPlaceFromText2(ctx context.Context, query string, token maps.PlaceAutocompleteSessionToken) (*Place2, error) {
	const op = "gmaps.Client2.findPlaceFromText2"

	resp, err := c.client.FindPlaceFromText(ctx, &maps.FindPlaceFromTextRequest{
		Input:        query,
		InputType:    maps.FindPlaceFromTextInputTypeTextQuery,
		Language:     c.cfg.Language,
		LocationBias: maps.FindPlaceFromTextLocationBiasCircular,
		Fields: []maps.PlaceSearchFieldMask{
			maps.PlaceSearchFieldMaskFormattedAddress,
			maps.PlaceSearchFieldMaskGeometryLocation,
			maps.PlaceSearchFieldMaskName,
			maps.PlaceSearchFieldMaskPlaceID,
			maps.PlaceSearchFieldMaskTypes,
		},
		LocationBiasCenter: &c.cfg.SearchCenter,
		LocationBiasRadius: int(c.cfg.SearchRadiusMeters),
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceNotFound2)
	}

	sawOutsideCoverage := false

	for _, candidate := range resp.Candidates {
		placeID := strings.TrimSpace(candidate.PlaceID)
		if placeID == "" {
			place := c.placeFromSearchResult2(candidate)
			if !place.InCoverage {
				sawOutsideCoverage = true
				continue
			}
			return place, nil
		}

		place, err := c.PlaceDetails2(ctx, placeID, token)
		switch {
		case err == nil:
			return place, nil
		case errors.Is(err, ErrPlaceOutOfCoverage2):
			sawOutsideCoverage = true
		case errors.Is(err, ErrPlaceNotFound2):
		default:
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	if sawOutsideCoverage {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceOutOfCoverage2)
	}

	return nil, fmt.Errorf("%s: %w", op, ErrPlaceNotFound2)
}

func (c *Client2) geocode2(ctx context.Context, query string) (*Place2, error) {
	const op = "gmaps.Client2.geocode2"

	results, err := c.client.Geocode(ctx, &maps.GeocodingRequest{
		Address:  query,
		Bounds:   &c.cfg.SearchBounds,
		Region:   c.cfg.Region,
		Language: c.cfg.Language,
		Components: map[maps.Component]string{
			maps.ComponentCountry: c.cfg.CountryCode,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result := c.pickBestGeocodingResult2(results)
	if result == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceNotFound2)
	}

	place := c.placeFromGeocode2(*result, "geocode")
	if !place.InCoverage {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceOutOfCoverage2)
	}

	return place, nil
}

func (c *Client2) validateClient2(op string) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("%s: %w", op, ErrClient2Nil)
	}

	return nil
}

func normalizeClient2Config(cfg Client2Config) Client2Config {
	def := DefaultClient2Config()

	if strings.TrimSpace(cfg.Language) == "" {
		cfg.Language = def.Language
	}
	if strings.TrimSpace(cfg.Region) == "" {
		cfg.Region = def.Region
	}
	if strings.TrimSpace(cfg.CountryCode) == "" {
		cfg.CountryCode = def.CountryCode
	}
	if cfg.SearchCenter.Lat == 0 && cfg.SearchCenter.Lng == 0 {
		cfg.SearchCenter = def.SearchCenter
	}
	if cfg.SearchRadiusMeters == 0 {
		cfg.SearchRadiusMeters = def.SearchRadiusMeters
	}
	if cfg.AutocompleteLimit <= 0 {
		cfg.AutocompleteLimit = def.AutocompleteLimit
	}
	if cfg.SearchBounds.NorthEast.Lat == 0 && cfg.SearchBounds.NorthEast.Lng == 0 &&
		cfg.SearchBounds.SouthWest.Lat == 0 && cfg.SearchBounds.SouthWest.Lng == 0 {
		cfg.SearchBounds = def.SearchBounds
	}

	cfg.Language = strings.TrimSpace(cfg.Language)
	cfg.Region = strings.TrimSpace(cfg.Region)
	cfg.CountryCode = strings.ToUpper(strings.TrimSpace(cfg.CountryCode))

	return cfg
}

func normalizeQuery2(input string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(input)), " ")
}

func isValidLatLng2(lat, lng float64) bool {
	if lat < -90 || lat > 90 {
		return false
	}
	if lng < -180 || lng > 180 {
		return false
	}
	return !(lat == 0 && lng == 0)
}

func cloneStrings2(items []string) []string {
	if len(items) == 0 {
		return nil
	}

	dst := make([]string, len(items))
	copy(dst, items)
	return dst
}

func (c *Client2) placeFromDetails2(result maps.PlaceDetailsResult) *Place2 {
	place := &Place2{
		Ref:         strings.TrimSpace(result.PlaceID),
		Source:      "place_details",
		Name:        normalizeTitleOrKeep2(result.Name),
		Address:     strings.TrimSpace(result.FormattedAddress),
		Lat:         result.Geometry.Location.Lat,
		Lng:         result.Geometry.Location.Lng,
		ResultTypes: cloneStrings2(result.Types),
	}

	c.fillAddressComponents2(place, result.AddressComponents)
	c.finalizePlace2(place)

	return place
}

func (c *Client2) placeFromSearchResult2(result maps.PlacesSearchResult) *Place2 {
	place := &Place2{
		Ref:         strings.TrimSpace(result.PlaceID),
		Source:      "find_place_from_text",
		Name:        normalizeTitleOrKeep2(result.Name),
		Address:     strings.TrimSpace(result.FormattedAddress),
		Lat:         result.Geometry.Location.Lat,
		Lng:         result.Geometry.Location.Lng,
		ResultTypes: cloneStrings2(result.Types),
	}

	c.finalizePlace2(place)

	return place
}

func (c *Client2) placeFromGeocode2(result maps.GeocodingResult, source string) *Place2 {
	place := &Place2{
		Ref:              strings.TrimSpace(result.PlaceID),
		Source:           source,
		Address:          strings.TrimSpace(result.FormattedAddress),
		Lat:              result.Geometry.Location.Lat,
		Lng:              result.Geometry.Location.Lng,
		ResultTypes:      cloneStrings2(result.Types),
		LocationType:     strings.TrimSpace(result.Geometry.LocationType),
		PartialMatch:     result.PartialMatch,
		PlusCodeGlobal:   strings.TrimSpace(result.PlusCode.GlobalCode),
		PlusCodeCompound: strings.TrimSpace(result.PlusCode.CompoundCode),
	}

	c.fillAddressComponents2(place, result.AddressComponents)
	c.finalizePlace2(place)
	place.Approximate = place.Approximate || result.PartialMatch

	return place
}

func (c *Client2) fillAddressComponents2(place *Place2, components []maps.AddressComponent) {
	for _, component := range components {
		for _, t := range component.Types {
			switch t {
			case "route":
				place.Route = normalizeStreet2(component.ShortName)
			case "street_number":
				place.StreetNumber = normalizeStreet2(component.ShortName)
			case "premise":
				place.Premise = normalizeTitleOrKeep2(component.LongName)
			case "subpremise":
				place.Subpremise = normalizeTitleOrKeep2(component.LongName)
			case "neighborhood", "sublocality_level_2":
				place.Neighborhood = normalizeTitleOrKeep2(component.LongName)
			case "sublocality":
				place.Sublocality = normalizeTitleOrKeep2(component.LongName)
			case "sublocality_level_1", "administrative_area_level_3":
				place.Commune = normalizeTitleOrKeep2(component.LongName)
			case "locality":
				place.Locality = normalizeTitleOrKeep2(component.LongName)
			case "administrative_area_level_1":
				place.AdministrativeL1 = normalizeTitleOrKeep2(component.LongName)
			case "administrative_area_level_2":
				place.AdministrativeL2 = normalizeTitleOrKeep2(component.LongName)
			case "country":
				place.Country = normalizeTitleOrKeep2(component.LongName)
				place.CountryCode = strings.ToUpper(strings.TrimSpace(component.ShortName))
			case "postal_code":
				place.PostalCode = strings.TrimSpace(component.LongName)
			}
		}
	}
}

func (c *Client2) finalizePlace2(place *Place2) {
	place.Route = strings.TrimSpace(place.Route)
	place.StreetNumber = strings.TrimSpace(place.StreetNumber)
	place.Address = strings.TrimSpace(place.Address)

	switch {
	case place.Route != "" && place.StreetNumber != "":
		place.AddressLine = place.Route + " # " + place.StreetNumber
	case place.Route != "":
		place.AddressLine = place.Route
	case place.Premise != "":
		place.AddressLine = place.Premise
	case place.Address != "":
		place.AddressLine = place.Address
	}

	if place.Name == "" {
		place.Name = place.Premise
	}

	if place.Locality == "" && place.AdministrativeL2 != "" {
		place.Locality = place.AdministrativeL2
	}

	if place.Commune == "" && place.Sublocality != "" {
		place.Commune = place.Sublocality
	}

	place.InCoverage = c.inCoverage2(place.Lat, place.Lng)
	place.Approximate = place.Approximate || isApproximateLocationType2(place.LocationType)
}

func (c *Client2) inCoverage2(lat, lng float64) bool {
	return lat >= c.cfg.SearchBounds.SouthWest.Lat &&
		lat <= c.cfg.SearchBounds.NorthEast.Lat &&
		lng >= c.cfg.SearchBounds.SouthWest.Lng &&
		lng <= c.cfg.SearchBounds.NorthEast.Lng
}

func (c *Client2) pickBestGeocodingResult2(results []maps.GeocodingResult) *maps.GeocodingResult {
	if len(results) == 0 {
		return nil
	}

	bestIndex := -1
	bestScore := -1

	for i := range results {
		score := c.geocodeScore2(results[i])
		if score > bestScore {
			bestScore = score
			bestIndex = i
		}
	}

	if bestIndex < 0 {
		return nil
	}

	return &results[bestIndex]
}

func (c *Client2) geocodeScore2(result maps.GeocodingResult) int {
	score := 0

	for _, t := range result.Types {
		switch t {
		case "street_address":
			score += 100
		case "premise":
			score += 95
		case "subpremise":
			score += 90
		case "route":
			score += 70
		case "neighborhood":
			score += 55
		case "sublocality", "sublocality_level_1":
			score += 45
		case "locality":
			score += 35
		}
	}

	switch strings.ToUpper(strings.TrimSpace(result.Geometry.LocationType)) {
	case "ROOFTOP":
		score += 20
	case "RANGE_INTERPOLATED":
		score += 10
	case "GEOMETRIC_CENTER":
		score += 5
	}

	if !result.PartialMatch {
		score += 10
	}

	if c.inCoverage2(result.Geometry.Location.Lat, result.Geometry.Location.Lng) {
		score += 25
	}

	return score
}

func normalizeStreet2(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	return normalizex.NormalizarStreet(value)
}

func normalizeTitleOrKeep2(value string) string {
	value = normalizeQuery2(value)
	if value == "" {
		return ""
	}

	return normalizex.NormalizeTitle(value)
}

func isApproximateLocationType2(value string) bool {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "APPROXIMATE", "GEOMETRIC_CENTER":
		return true
	default:
		return false
	}
}
