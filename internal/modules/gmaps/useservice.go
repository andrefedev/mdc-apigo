package gmaps

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"googlemaps.github.io/maps"
)

func (c *Client) Autocomplete(ctx context.Context, input string, token maps.PlaceAutocompleteSessionToken) ([]*Prediction, error) {
	const op = "gmaps.Client.Autocomplete"

	query := normalizeQuery(input)
	if query == "" {
		return nil, fmt.Errorf("%s: %w", op, ErrQueryRequired)
	}

	resp, err := c.client.PlaceAutocomplete(ctx, &maps.PlaceAutocompleteRequest{
		Input:        query,
		Origin:       &c.config.SearchCenter,
		Radius:       c.config.SearchRadiusMeters,
		Location:     &c.config.SearchCenter,
		Language:     c.config.Language,
		StrictBounds: c.config.StrictBounds,
		SessionToken: token,
		Components: map[maps.Component][]string{
			maps.ComponentCountry: {c.config.CountryCode},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	limit := c.config.AutocompleteLimit
	if limit <= 0 || limit > len(resp.Predictions) {
		limit = len(resp.Predictions)
	}

	results := make([]*Prediction, 0, limit)
	for i := 0; i < limit; i++ {
		p := resp.Predictions[i]
		results = append(results, &Prediction{
			Ref:            strings.TrimSpace(p.PlaceID),
			Desc:           strings.TrimSpace(p.Description),
			Title:          strings.TrimSpace(p.StructuredFormatting.MainText),
			Subtitle:       strings.TrimSpace(p.StructuredFormatting.SecondaryText),
			DistanceMeters: int32(p.DistanceMeters),
			// Types: cloneStrings2(p.Types),
		})
	}

	return results, nil
}

func (c *Client) ResolveText(ctx context.Context, input string, token maps.PlaceAutocompleteSessionToken) (*Place, error) {
	const op = "gmaps.Client.ResolveText"

	query := normalizeQuery(input)
	if query == "" {
		return nil, fmt.Errorf("%s: %w", op, ErrQueryRequired)
	}

	sawOutsideCoverage := false

	predictions, err := c.Autocomplete(ctx, query, token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for _, prediction := range predictions {
		place, err := c.PlaceDetails(ctx, prediction.Ref, token)
		switch {
		case err == nil:
			return place, nil
		case errors.Is(err, ErrPlaceOutOfCoverage):
			sawOutsideCoverage = true
		case errors.Is(err, ErrPlaceNotFound):
		default:
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	place, err := c.findPlaceFromText(ctx, query, token)
	switch {
	case err == nil:
		return place, nil
	case errors.Is(err, ErrPlaceOutOfCoverage):
		sawOutsideCoverage = true
	case errors.Is(err, ErrPlaceNotFound):
	default:
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	place, err = c.geocode(ctx, query)
	switch {
	case err == nil:
		return place, nil
	case errors.Is(err, ErrPlaceOutOfCoverage):
		sawOutsideCoverage = true
	case errors.Is(err, ErrPlaceNotFound):
	default:
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if sawOutsideCoverage {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceOutOfCoverage)
	}

	return nil, fmt.Errorf("%s: %w", op, ErrPlaceNotFound)
}

func (c *Client) PlaceDetails(ctx context.Context, placeID string, token maps.PlaceAutocompleteSessionToken) (*Place, error) {
	const op = "gmaps.Client.PlaceDetails"

	placeID = strings.TrimSpace(placeID)
	if placeID == "" {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceRefRequired)
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
		Language:     c.config.Language,
		Region:       c.config.Region,
		Fields:       fields,
		SessionToken: token,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	place := c.placeFromDetails(result)
	if !place.InCoverage {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceOutOfCoverage)
	}

	return place, nil
}

func (c *Client) ReverseGeocode(ctx context.Context, lat, lng float64) (*Place, error) {
	const op = "gmaps.Client.ReverseGeocode"

	if !isValidLatLng(lat, lng) {
		return nil, fmt.Errorf("%s: %w", op, ErrCoordinatesInvalid)
	}

	latLng := &maps.LatLng{Lat: lat, Lng: lng}

	results, err := c.client.ReverseGeocode(ctx, &maps.GeocodingRequest{
		LatLng:     latLng,
		Region:     c.config.Region,
		Language:   c.config.Language,
		ResultType: []string{"street_address", "premise", "subpremise"},
		LocationType: []maps.GeocodeAccuracy{
			maps.GeocodeAccuracyRooftop,
			maps.GeocodeAccuracyRangeInterpolated,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(results) == 0 {
		results, err = c.client.ReverseGeocode(ctx, &maps.GeocodingRequest{
			LatLng:     latLng,
			Region:     c.config.Region,
			Language:   c.config.Language,
			ResultType: []string{"route", "neighborhood", "sublocality", "locality"},
		})
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	result := c.pickBestGeocodingResult(results)
	if result == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceNotFound)
	}

	place := c.placeFromGeocode(*result)
	if !place.InCoverage {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceOutOfCoverage)
	}

	return place, nil
}

func (c *Client) findPlaceFromText(ctx context.Context, query string, token maps.PlaceAutocompleteSessionToken) (*Place, error) {
	const op = "gmaps.Client.findPlaceFromText"

	resp, err := c.client.FindPlaceFromText(ctx, &maps.FindPlaceFromTextRequest{
		Input:        query,
		InputType:    maps.FindPlaceFromTextInputTypeTextQuery,
		Language:     c.config.Language,
		LocationBias: maps.FindPlaceFromTextLocationBiasCircular,
		Fields: []maps.PlaceSearchFieldMask{
			maps.PlaceSearchFieldMaskFormattedAddress,
			maps.PlaceSearchFieldMaskGeometryLocation,
			maps.PlaceSearchFieldMaskName,
			maps.PlaceSearchFieldMaskPlaceID,
			maps.PlaceSearchFieldMaskTypes,
		},
		LocationBiasCenter: &c.config.SearchCenter,
		LocationBiasRadius: int(c.config.SearchRadiusMeters),
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceNotFound)
	}

	sawOutsideCoverage := false

	for _, candidate := range resp.Candidates {
		placeID := strings.TrimSpace(candidate.PlaceID)
		if placeID == "" {
			place := c.placeFromSearchResult(candidate)
			if !place.InCoverage {
				sawOutsideCoverage = true
				continue
			}
			return place, nil
		}

		place, err := c.PlaceDetails(ctx, placeID, token)
		switch {
		case err == nil:
			return place, nil
		case errors.Is(err, ErrPlaceOutOfCoverage):
			sawOutsideCoverage = true
		case errors.Is(err, ErrPlaceNotFound):
		default:
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	if sawOutsideCoverage {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceOutOfCoverage)
	}

	return nil, fmt.Errorf("%s: %w", op, ErrPlaceNotFound)
}

func (c *Client) geocode(ctx context.Context, query string) (*Place, error) {
	const op = "gmaps.Client.geocode"

	results, err := c.client.Geocode(ctx, &maps.GeocodingRequest{
		Address:  query,
		Region:   c.config.Region,
		Bounds:   &c.config.SearchBounds,
		Language: c.config.Language,
		Components: map[maps.Component]string{
			maps.ComponentCountry: c.config.CountryCode,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result := c.pickBestGeocodingResult(results)
	if result == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceNotFound)
	}

	place := c.placeFromGeocode(*result)
	if !place.InCoverage {
		return nil, fmt.Errorf("%s: %w", op, ErrPlaceOutOfCoverage)
	}

	return place, nil
}

func (c *Client) placeFromDetails(result maps.PlaceDetailsResult) *Place {
	place := &Place{
		Ref:     strings.TrimSpace(result.PlaceID),
		Name:    normalizeQuery(result.Name),
		Address: strings.TrimSpace(result.FormattedAddress),
		Lat:     result.Geometry.Location.Lat,
		Lng:     result.Geometry.Location.Lng,
	}

	c.fillAddressComponents(place, result.AddressComponents)
	c.finalizePlace(place)

	return place
}

func (c *Client) placeFromSearchResult(result maps.PlacesSearchResult) *Place {
	place := &Place{
		Ref:     strings.TrimSpace(result.PlaceID),
		Name:    normalizeQuery(result.Name),
		Address: strings.TrimSpace(result.FormattedAddress),
		Lat:     result.Geometry.Location.Lat,
		Lng:     result.Geometry.Location.Lng,
	}

	c.finalizePlace(place)

	return place
}

func (c *Client) placeFromGeocode(result maps.GeocodingResult) *Place {
	place := &Place{
		Ref:     strings.TrimSpace(result.PlaceID),
		Address: strings.TrimSpace(result.FormattedAddress),
		Lat:     result.Geometry.Location.Lat,
		Lng:     result.Geometry.Location.Lng,
	}

	c.fillAddressComponents(place, result.AddressComponents)
	c.finalizePlace(place)
	place.Approximate = result.PartialMatch || isApproximateLocationType(result.Geometry.LocationType)

	return place
}

func (c *Client) fillAddressComponents(place *Place, components []maps.AddressComponent) {
	var premise string

	for _, component := range components {
		for _, t := range component.Types {
			switch t {
			case "route":
				place.Route = normalizeStreet(component.ShortName)
			case "street_number":
				place.Street = normalizeStreet(component.ShortName)
			case "premise":
				premise = normalizeQuery(component.LongName)
			case "neighborhood", "sublocality_level_2":
				place.Neighb = normalizeTitleOrKeep(component.LongName)
			case "sublocality":
				place.Sublocal = normalizeTitleOrKeep(component.LongName)
			case "sublocality_level_1", "administrative_area_level_3":
				place.Cmna = normalizeTitleOrKeep(component.LongName)
			case "locality", "administrative_area_level_2":
				place.Locality = normalizeTitleOrKeep(component.LongName)
			}
		}
	}

	if place.Name == "" {
		place.Name = premise
	}
}

func (c *Client) finalizePlace(place *Place) {
	place.Route = strings.TrimSpace(place.Route)
	place.Street = strings.TrimSpace(place.Street)
	place.Address = strings.TrimSpace(place.Address)

	place.InCoverage = c.inCoverage(place.Lat, place.Lng)
}

func (c *Client) inCoverage(lat, lng float64) bool {
	return lat >= c.config.SearchBounds.SouthWest.Lat &&
		lat <= c.config.SearchBounds.NorthEast.Lat &&
		lng >= c.config.SearchBounds.SouthWest.Lng &&
		lng <= c.config.SearchBounds.NorthEast.Lng
}

func (c *Client) pickBestGeocodingResult(results []maps.GeocodingResult) *maps.GeocodingResult {
	if len(results) == 0 {
		return nil
	}

	bestIndex := -1
	bestScore := -1

	for i := range results {
		score := c.geocodeScore(results[i])
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

func (c *Client) geocodeScore(result maps.GeocodingResult) int {
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

	if c.inCoverage(result.Geometry.Location.Lat, result.Geometry.Location.Lng) {
		score += 25
	}

	return score
}
