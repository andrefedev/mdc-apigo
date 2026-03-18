package messages

// NewTemplateMessageReq maps template input into a Graph API payload.
func NewTemplateMessageReq(input *SendTemplateMessage) TemplateMessageData {
	// parameters := make([]TemplateParameter, 0, len(input.BodyText))

	//for _, item := range input.BodyText {
	//	text := item
	//	parameters = append(parameters, TemplateParameter{
	//		Type: TemplateParameterTypeText,
	//		Text: &text,
	//	})
	//}

	//components := make([]TemplateComponent, 0, 1)
	//if len(parameters) > 0 {
	//	components = append(components, TemplateComponent{
	//		Type:       "body",
	//		Parameters: parameters,
	//	})
	//}

	templ := input.Template
	templ.Language.Code = "es_CO"

	return TemplateMessageData{
		To:               "57" + input.To,
		Type:             TypeTemplate,
		Template:         templ,
		RecipientType:    RecipientTypeIndividual,
		MessagingProduct: MessagingProductWhatsApp,
		// Context:          input.Context,
		//Template: &TemplateContent{
		//	Name:       input.Name,
		//	Language:   input.Language,
		//	Components: components,
		//},
	}
}
