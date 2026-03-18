package messages

type TemplateMessageRequest struct {
	To       string        `json:"to"`
	Type     string        `json:"type"` // "template"
	Template *TemplContent `json:"template"`
}
