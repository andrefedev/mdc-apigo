package gmaps

// ###############
// # GOOGLE MAPS #
// ###############

type Place struct {
	Ref      string
	Lat      float64
	Lng      float64
	Name     string
	Cmna     string
	Route    string
	Street   string
	Neighb   string
	Address  string
	Locality string
	Sublocal string
}

//func (r *Place) ToProto() *v1.Place {
//	return &v1.Place{
//		Ref:      r.Ref,
//		Lat:      r.Lat,
//		Lng:      r.Lng,
//		Name:     r.Name,
//		Cmna:     r.Cmna,
//		Route:    r.Route,
//		Street:   r.Street,
//		Neighb:   r.Neighb,
//		Address:  r.Address,
//		Locality: r.Locality,
//		Sublocal: r.Sublocal,
//	}
//}

type Prediction struct {
	Ref      string `json:"ref"`
	Desc     string `json:"desc"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

//func (r *Prediction) ToProto() *v1.Prediction {
//	return &v1.Prediction{
//		Ref:      r.Ref,
//		Desc:     r.Desc,
//		Title:    r.Title,
//		Subtitle: r.Subtitle,
//	}
//}

type Geocoding struct {
	Ref              string  `json:"ref"`
	Route            string  `json:"route"`
	Street           string  `json:"street"`
	Locality         string  `json:"locality"`
	Sublocal         string  `json:"sublocal"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	PartialMatch     bool    `json:"partialMatch"`
	FormattedAddress string  `json:"formattedAddress"`
}
