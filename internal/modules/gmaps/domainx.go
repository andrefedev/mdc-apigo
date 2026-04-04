package gmaps

import (
	"apigo/internal/platforms/validatex/normalizex"
	v1 "apigo/protobuf/gen/v1"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// PLACE_AUTOCOMPLETE_DATA__

type PlaceAutocompleteData struct {
	Query string
}

func NewPlaceAutocompleteData(input *PlaceAutocompleteInput) *PlaceAutocompleteData {
	if input == nil {
		return &PlaceAutocompleteData{}
	}

	return &PlaceAutocompleteData{
		Query: input.Query,
	}
}

func (d *PlaceAutocompleteData) Validate() error {
	const op = "Mapx.PlaceAutocompleteData.Validate"

	// Normalize
	d.Query = strings.TrimSpace(d.Query)
	d.Query = normalizex.NormalizeName(d.Query)

	// Validation
	if d.Query == "" {
		return fmt.Errorf("%s: %w", op, WrapQueryRequired(nil))
	}

	return nil
}

// PLACE_DETAIL_DATA__

type PlaceDetailData struct {
	Ref   string
	Token string
}

func NewPlaceDetailData(input *PlaceDetailInput) *PlaceDetailData {
	if input == nil {
		return &PlaceDetailData{}
	}

	return &PlaceDetailData{
		Ref:   input.Ref,
		Token: input.Token,
	}
}

func (d *PlaceDetailData) Validate() error {
	const op = "Mapx.PlaceDetailData.Validate"

	// Normalize
	d.Ref = strings.TrimSpace(d.Ref)
	d.Token = strings.TrimSpace(d.Token)

	// Validation
	if d.Ref == "" {
		return fmt.Errorf("%s: %w", op, WrapPlaceRefRequired(nil))
	}
	if d.Token == "" {
		return fmt.Errorf("%s: %w", op, WrapPlaceTokenRequired(nil))
	}

	return nil
}

// REVERSE_GEOCODE_DATA__

type ReverseGeocodeData struct {
	Lat float64
	Lng float64
}

func NewReverseGeocodeData(input *ReverseGeocodeInput) *ReverseGeocodeData {
	if input == nil {
		return &ReverseGeocodeData{}
	}

	return &ReverseGeocodeData{
		Lat: input.Lat,
		Lng: input.Lng,
	}
}

func (d *ReverseGeocodeData) Validate() error {
	const op = "Mapx.ReverseGeocodeData.Validate"

	// Validation
	if d.Lat < -90 || d.Lat > 90 {
		return fmt.Errorf("%s: %w", op, WrapCoordinatesInvalid(nil))
	}
	if d.Lng < -180 || d.Lng > 180 {
		return fmt.Errorf("%s: %w", op, WrapCoordinatesInvalid(nil))
	}
	if d.Lat == 0 && d.Lng == 0 {
		return fmt.Errorf("%s: %w", op, WrapCoordinatesInvalid(nil))
	}

	return nil
}

// ##################

// PLACE_AUTOCOMPLETE_INPUT__

type PlaceAutocompleteInput struct {
	Query string
}

func NewPlaceAutocompleteInput(req *v1.PlaceAutocompleteReq) *PlaceAutocompleteInput {
	if req == nil {
		return &PlaceAutocompleteInput{}
	}

	return &PlaceAutocompleteInput{
		Query: req.GetQuery(),
	}
}

func (r *PlaceAutocompleteInput) Validate() error {
	const op = "Mapx.PlaceAutocompleteInput.Validate"

	// Normalize
	r.Query = strings.TrimSpace(r.Query)
	r.Query = normalizex.NormalizeName(r.Query)

	// Validation
	if r.Query == "" {
		return fmt.Errorf("%s: %w", op, WrapQueryRequired(nil))
	}

	return nil
}

// PLACE_DETAIL_INPUT__

type PlaceDetailInput struct {
	Ref   string
	Token string
}

func NewPlaceDetailInput(req *v1.PlaceDetailReq) *PlaceDetailInput {
	if req == nil {
		return &PlaceDetailInput{}
	}

	return &PlaceDetailInput{
		Ref:   strings.TrimSpace(req.GetRef()),
		Token: strings.TrimSpace(req.GetToken()),
	}
}

func (r *PlaceDetailInput) Validate() error {
	const op = "Mapx.PlaceDetailInput.Validate"

	// Normalize
	r.Ref = strings.TrimSpace(r.Ref)
	r.Token = strings.TrimSpace(r.Token)

	// Validation
	if r.Ref == "" {
		return fmt.Errorf("%s: %w", op, WrapPlaceRefRequired(nil))
	}
	if r.Token == "" {
		return fmt.Errorf("%s: %w", op, WrapPlaceTokenRequired(nil))
	}
	if err := uuid.Validate(r.Token); err != nil {
		return fmt.Errorf("%s: %w", op, WrapPlaceTokenInvalid(err))
	}

	return nil
}

// REVERSE_GEOCODE_INPUT__

type ReverseGeocodeInput struct {
	Lat float64
	Lng float64
}

func NewReverseGeocodeInput(req *v1.ReverseGeocodeReq) *ReverseGeocodeInput {
	if req == nil {
		return &ReverseGeocodeInput{}
	}

	return &ReverseGeocodeInput{
		Lat: req.GetLat(),
		Lng: req.GetLng(),
	}
}

func (r *ReverseGeocodeInput) Validate() error {
	const op = "Mapx.ReverseGeocodeInput.Validate"

	if r.Lat < -90 || r.Lat > 90 {
		return fmt.Errorf("%s: %w", op, WrapCoordinatesInvalid(nil))
	}
	if r.Lng < -180 || r.Lng > 180 {
		return fmt.Errorf("%s: %w", op, WrapCoordinatesInvalid(nil))
	}
	if r.Lat == 0 && r.Lng == 0 {
		return fmt.Errorf("%s: %w", op, WrapCoordinatesInvalid(nil))
	}

	return nil
}
