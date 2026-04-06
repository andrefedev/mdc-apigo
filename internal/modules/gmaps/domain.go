package gmaps

import v1 "apigo/protobuf/gen/v1"

type Place struct {
	Ref         string  `json:"ref"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Name        string  `json:"name,omitempty"`
	Cmna        string  `json:"cmna,omitempty"`
	Route       string  `json:"route,omitempty"`
	Street      string  `json:"street,omitempty"`
	Neighb      string  `json:"neighb,omitempty"`
	Address     string  `json:"address,omitempty"`
	Locality    string  `json:"locality,omitempty"`
	Sublocal    string  `json:"sublocal,omitempty"`
	InCoverage  bool    `json:"inCoverage"`
	Approximate bool    `json:"approximate"`
}

func (p *Place) ToProto() *v1.Place {
	if p == nil {
		return nil
	}

	return &v1.Place{
		Ref:         p.Ref,
		Lat:         p.Lat,
		Lng:         p.Lng,
		Name:        p.Name,
		Cmna:        p.Cmna,
		Route:       p.Route,
		Street:      p.Street,
		Neighb:      p.Neighb,
		Address:     p.Address,
		Locality:    p.Locality,
		Sublocal:    p.Sublocal,
		InCoverage:  p.InCoverage,
		Approximate: p.Approximate,
	}
}

type Prediction struct {
	Ref            string `json:"ref"`
	Desc           string `json:"desc"`
	Title          string `json:"title"`
	Subtitle       string `json:"subtitle"`
	DistanceMeters int32  `json:"distanceMeters,omitempty"`
}

func (p *Prediction) ToProto() *v1.Prediction {
	if p == nil {
		return nil
	}

	return &v1.Prediction{
		Ref:            p.Ref,
		Desc:           p.Desc,
		Title:          p.Title,
		Subtitle:       p.Subtitle,
		DistanceMeters: p.DistanceMeters,
	}
}
