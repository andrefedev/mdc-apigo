package messages

//// sendMessageResponse2 mirrors the WhatsApp Graph API success payload for sends.
//type sendMessageResponse2 struct {
//	MessagingProduct string                `json:"messaging_product"`
//	Contacts         []sendMessageContact2 `json:"contacts"`
//	Messages         []sendMessageMessage2 `json:"messages"`
//}
//
//type sendMessageContact2 struct {
//	Input string `json:"input"`
//	WaID  string `json:"wa_id"`
//}
//
//type sendMessageMessage2 struct {
//	ID            string `json:"id"`
//	MessageStatus string `json:"message_status,omitempty"`
//}
//
//// SendTemplateResult2 is the normalized result exposed by Service2.
//type SendTemplateResult2 struct {
//	MessageID string
//	Status    string
//	Input     string
//	WaID      string
//}
