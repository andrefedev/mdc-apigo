package messages

type TemplLang struct {
	Code    string `json:"code"`
	Policty string `json:"policy,omitempty"`
}

type TemplComp struct {
	Type       string       `json:"type"`
	Index      *int         `json:"index,omitempty"`
	SubType    *string      `json:"sub_type,omitempty"`
	Parameters []TemplParam `json:"parameters,omitempty"`
}

type TemplParam struct {
	Type     string                  `json:"type"`
	Text     *string                 `json:"text,omitempty"`
	Image    *TemplMediaObject       `json:"image,omitempty"`
	Video    *TemplMediaObject       `json:"video,omitempty"`
	Payload  *string                 `json:"payload,omitempty"`
	Document *TemplMediaObject       `json:"document,omitempty"`
	Currency *TemplCurrencyParameter `json:"currency,omitempty"`
	DateTime *TemplDateTimeParameter `json:"date_time,omitempty"`
}

// TemplContent contains a template send payload.
type TemplContent struct {
	Name       string      `json:"name"`
	Language   TemplLang   `json:"language"`
	Components []TemplComp `json:"components,omitempty"`
}

type TemplateMessageData struct {
	To               string        `json:"to"`
	Type             string        `json:"type"` // "template"
	Template         *TemplContent `json:"template"`
	RecipientType    string        `json:"recipient_type,omitempty"` // "individual"
	MessagingProduct string        `json:"messaging_product"`        // "whatsapp"
}

type SendTemplateMessage struct {
	To       string        `json:"to"`
	Type     string        `json:"type"` // "template"
	Template *TemplContent `json:"template"`
	// RecipientType string `json:"recipient_type,omitempty"` // "individual"
	// MessagingProduct string `json:"messaging_product"` // "WhatsApp"
}
