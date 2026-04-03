package gmaps

type Place2 struct {
	Ref              string   `json:"ref"`
	Query            string   `json:"query,omitempty"`
	Source           string   `json:"source"`
	Name             string   `json:"name,omitempty"`
	Address          string   `json:"address"`
	AddressLine      string   `json:"addressLine,omitempty"`
	Route            string   `json:"route,omitempty"`
	StreetNumber     string   `json:"streetNumber,omitempty"`
	Premise          string   `json:"premise,omitempty"`
	Subpremise       string   `json:"subpremise,omitempty"`
	Neighborhood     string   `json:"neighborhood,omitempty"`
	Sublocality      string   `json:"sublocality,omitempty"`
	Commune          string   `json:"commune,omitempty"`
	Locality         string   `json:"locality,omitempty"`
	AdministrativeL1 string   `json:"administrativeL1,omitempty"`
	AdministrativeL2 string   `json:"administrativeL2,omitempty"`
	Country          string   `json:"country,omitempty"`
	CountryCode      string   `json:"countryCode,omitempty"`
	PostalCode       string   `json:"postalCode,omitempty"`
	PlusCodeGlobal   string   `json:"plusCodeGlobal,omitempty"`
	PlusCodeCompound string   `json:"plusCodeCompound,omitempty"`
	Lat              float64  `json:"lat"`
	Lng              float64  `json:"lng"`
	ResultTypes      []string `json:"resultTypes,omitempty"`
	LocationType     string   `json:"locationType,omitempty"`
	PartialMatch     bool     `json:"partialMatch"`
	Approximate      bool     `json:"approximate"`
	InCoverage       bool     `json:"inCoverage"`
}

type Prediction2 struct {
	Ref            string   `json:"ref"`
	Desc           string   `json:"desc"`
	Title          string   `json:"title"`
	Subtitle       string   `json:"subtitle"`
	DistanceMeters int      `json:"distanceMeters,omitempty"`
	Types          []string `json:"types,omitempty"`
}
