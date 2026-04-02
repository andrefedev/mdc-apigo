package gmaps

import (
	"apigo"
	"context"
	"fmt"
	"log"

	"googlemaps.github.io/maps"
)

type MapsVendor struct {
	client *maps.Client
}

// defaultValue
var latLng = &maps.LatLng{
	Lat: 6.249623173934768,  //6.249639919536681, 6.249623173934768, -75.58688348047991
	Lng: -75.58688348047991, //-75.58688791329203,
}

var latLng2 = &maps.LatLng{
	Lat: 6.266733118358222,
	Lng: -75.58688348047991,
}

func NewMapsVendor(client *maps.Client) *MapsVendor {
	return &MapsVendor{client}
}

// PlaceDetails función
func (r MapsVendor) PlaceDetails(ctx context.Context, placeRef string, token maps.PlaceAutocompleteSessionToken) (*domain.Place, error) {
	fields := []maps.PlaceDetailsFieldMask{
		maps.PlaceDetailsFieldMaskName,
		maps.PlaceDetailsFieldMaskPlaceID,
		maps.PlaceDetailsFieldMaskADRAddress,
		maps.PlaceDetailsFieldMaskFormattedAddress,
		maps.PlaceDetailsFieldMaskAddressComponent,
		maps.PlaceDetailsFieldMaskGeometryLocation,
	}

	request := &maps.PlaceDetailsRequest{
		PlaceID:      placeRef,
		Language:     "es",
		Region:       "co",
		Fields:       fields,
		SessionToken: token,
	}

	pd, err := r.client.PlaceDetails(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("MapsVendor.PlaceDetails: [%w]", err)
	}

	result := &domain.Place{
		Ref:     pd.PlaceID,
		Lat:     pd.Geometry.Location.Lat,
		Lng:     pd.Geometry.Location.Lng,
		Name:    apigo.NormalizeTitle(pd.Name),
		Address: pd.FormattedAddress,
	}

	for _, comp := range pd.AddressComponents {
		for _, t := range comp.Types {
			switch t {
			case "sublocality_level_1",
				"administrative_area_level_3": // cmna
				result.Cmna = comp.LongName
			case "route":
				result.Route = comp.ShortName
			case "street_number":
				result.Street = comp.ShortName
			case "locality",
				"administrative_area_level_2":
				result.Locality = comp.LongName
			case "neighborhood", "sublocality_level_2":
				result.Neighb = comp.LongName
			case "sublocality":
				result.Sublocal = comp.LongName
			}
		}

		fmt.Printf("PlaceDetails: %s -> %v\n", comp.LongName, comp.Types)
	}

	return result, nil
}

// ReverseGeocode función
func (r MapsVendor) ReverseGeocode(ctx context.Context, lat, lng float64) (*domain.Place, error) {
	req := &maps.GeocodingRequest{
		LatLng:       &maps.LatLng{Lat: lat, Lng: lng},
		Region:       "co",
		Language:     "es",
		ResultType:   []string{"street_address", "premise", "subpremise"}, // opcional: filtra a direcciones
		LocationType: []maps.GeocodeAccuracy{maps.GeocodeAccuracyRooftop, maps.GeocodeAccuracyRangeInterpolated},
	}

	res, err := r.client.ReverseGeocode(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("MapsVendor.ReverseGeocode: [%w]", err)
	}

	// empty
	if len(res) == 0 {
		return nil, nil
	}

	first := res[0]
	result := &domain.Place{
		Ref:     first.PlaceID,
		Lat:     first.Geometry.Location.Lat,
		Lng:     first.Geometry.Location.Lng,
		Address: first.FormattedAddress,
	}

	for _, comp := range first.AddressComponents {
		for _, t := range comp.Types {
			switch t {
			case "route":
				result.Route = apigo.NormalizarStreet(comp.ShortName)
			case "street_number":
				result.Street = apigo.NormalizarStreet(comp.ShortName)
			case "locality":
				result.Locality = apigo.NormalizeTitle(comp.LongName) // city or locality
			case "neighborhood":
				result.Neighb = apigo.NormalizeTitle(comp.LongName) // barrio neighborhood
			}
		}

		fmt.Printf("ReverseGeocode %s -> %v\n", comp.LongName, comp.Types)
	}

	return result, nil
}

// PlaceAutocomplete función
func (r MapsVendor) PlaceAutocomplete(ctx context.Context, input string, token maps.PlaceAutocompleteSessionToken) ([]*domain.Prediction, error) {
	req := &maps.PlaceAutocompleteRequest{
		Input:        input,
		Radius:       25000, // 50km
		Language:     "es-419",
		Location:     latLng,
		StrictBounds: true, // true restringe demasiado; puedes alternar
		SessionToken: token,
		Components: map[maps.Component][]string{
			maps.ComponentCountry: {"CO"},
		},
	}

	res, err := r.client.PlaceAutocomplete(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("MapsVedor.PlaceAutocomqplete: [client place autocomplete] [%w]", err)
	}

	predictions := make([]*domain.Prediction, 0)
	for _, p := range res.Predictions {

		var prediction domain.Prediction
		prediction.Ref = p.PlaceID
		prediction.Desc = p.Description
		prediction.Title = p.StructuredFormatting.MainText
		prediction.Subtitle = p.StructuredFormatting.SecondaryText
		predictions = append(predictions, &prediction)
	}

	return predictions, nil
}

func (r MapsVendor) ComputeDistanceMatrix(ctx context.Context, dests []maps.LatLng) error {
	origin := latLng2.String()
	destination := latLng.String()

	waypoints := make([]string, len(dests))
	for i, d := range dests {
		waypoints[i] = d.String()
	}

	req := &maps.DirectionsRequest{
		Optimize:    true,
		Origin:      origin,
		Destination: destination,
		Waypoints:   waypoints,
		Mode:        maps.TravelModeDriving,

		// Optimize:    true,
		// Units:       maps.UnitsMetric,
		// Language:      "es",
		// Region:        "co",
		DepartureTime: "now",
		// TrafficModel:  maps.TrafficModelBestGuess,
	}

	routes, _, err := r.client.Directions(ctx, req)
	if err != nil {
		return fmt.Errorf("directions error: %w", err)
	}

	if len(routes) == 0 {
		return fmt.Errorf("directions no devolvió rutas")
	}

	optimizedOrder := routes[0].WaypointOrder
	fmt.Println("Orden óptimo:", optimizedOrder)

	// Construimos el intent URL para Android
	var waypointParts []string
	for _, i := range optimizedOrder {
		w := waypoints[i]
		waypointParts = append(waypointParts, w)
	}

	url := fmt.Sprintf(
		"google.navigation:q=%s&waypoints=%s",
		destination,
		joinWaypoints(waypointParts),
	)

	log.Printf("navURI: %s", url)

	return nil
}

func joinWaypoints(points []string) string {
	return fmt.Sprintf("%s", joinStrings(points, "|"))
}

func joinStrings(list []string, sep string) string {
	result := ""
	for i, s := range list {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

//func buildGoogleMapsURL(, waypointOrder []int) string {
//	origin := fmt.Sprintf("%.6f,%.6f", depotLat, depotLng)
//	dest := origin
//
//	ordered := make([]string, len(orders))
//	for i, idx := range waypointOrder {
//		o := orders[idx]
//		ordered[i] = fmt.Sprintf("%.6f,%.6f", o.Location.Lat, o.Location.Lng)
//	}
//
//	waypointsStr := ""
//	for i, w := range ordered {
//		if i > 0 {
//			waypointsStr += "|"
//		}
//		waypointsStr += w
//	}
//
//	q := url.Values{}
//	q.Set("api", "1")
//	q.Set("origin", origin)
//	q.Set("destination", dest)
//	q.Set("travelmode", "driving")
//	if waypointsStr != "" {
//		q.Set("waypoints", waypointsStr)
//	}
//
//	return "https://www.google.com/maps/dir/?" + q.Encode()
//}
