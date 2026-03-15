package messages

const (
	TypeText        = "text"
	TypeTemplate    = "template"
	TypeDocument    = "document"
	TypeInteractive = "interactive"

	InteractiveTypeFlow        = "flow"
	InteractiveTypeProduct     = "product"
	InteractiveTypeProductList = "product_list"

	InteractiveHeaderTypeText     = "text"
	InteractiveHeaderTypeImage    = "image"
	InteractiveHeaderTypeVideo    = "video"
	InteractiveHeaderTypeDocument = "document"

	//TemplCategoryUtility        = "UTILITY"
	//TemplCategoryMarketing      = "MARKETING"
	//TemplCategoryAuthentication = "AUTHENTICATION"

	TemplParameterTypeDoc   = "document"
	TemplParameterTypeText  = "text"
	TemplParameterTypeImage = "image"
	TemplParameterTypeVideo = "video"

	FlowActionNavigate     = "navigate"
	FlowActionDataExchange = "data_exchange"

	RecipientTypeIndividual  = "individual"
	MessagingProductWhatsApp = "whatsapp"
)
