package messages

type TemplMediaObject struct {
	Ref      string `json:"id,omitempty"`
	Link     string `json:"link,omitempty"`
	Caption  string `json:"caption,omitempty"`
	Filename string `json:"filename,omitempty"`
}

type TemplCurrencyParameter struct {
	Code          string `json:"code"`
	Amount1000    int64  `json:"amount_1000"`
	FallbackValue string `json:"fallback_value"`
}

type TemplDateTimeParameter struct {
	FallbackValue string `json:"fallback_value"`
}
