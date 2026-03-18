package messages

// NewTemplateMessageReq maps template input into a Graph API payload.
func newTemplateMessage(input *TemplateMessageRequest) TemplateMessage {
	// parameters := make([]TemplateParameter, 0, len(input.BodyText))

	//components := make([]TemplateComponent, 0, 1)
	//if len(parameters) > 0 {
	//	components = append(components, TemplateComponent{
	//		Type:       "body",
	//		Parameters: parameters,
	//	})
	//}

	return TemplateMessage{
		To:               input.To,
		Type:             TypeTemplate,
		Template:         input.Template,
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
