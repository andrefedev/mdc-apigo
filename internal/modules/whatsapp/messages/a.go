package messages

//type Language struct {
//	Code string `json:"code"`
//}
//
//// TemplateLanguage identifies the language policy and code.
//type TemplateLanguage struct {
//	Code   string `json:"code"`
//	Policy string `json:"policy,omitempty"`
//}
//
//type Template struct {
//	Name     string `json:"name"`
//	Language struct {
//		Code string `json:"code"`
//	} `json:"language"`
//}
//
//type TemplateMessage struct {
//	To       string   `json:"to"`
//	Type     string   `json:"type"`
//	Product  string   `json:"messaging_product"`
//	Template Template `json:"template"`
//}
//
//// ##############
//
//type MessageService struct {
//	client *Client
//}
//
//func NewMessageService(client *Client) *MessageService {
//	return &MessageService{client: client}
//}
//
//// SendTemplate sirve para Autenticación (OTP) u Ofertas
//func (s *MessageService) SendTemplate(ctx context.Context, to, lang, template string) error {
//	if s == nil || s.client == nil {
//		return errors.New("whatsapp message service is not configured")
//	}
//
//	reqBody := TemplateMessage{
//		To:      to,
//		Type:    "template",
//		Product: "whatsapp",
//	}
//
//	reqBody.Template.Name = template
//	reqBody.Template.Language.Code = lang
//
//	bodyBytes, err := json.Marshal(reqBody)
//	if err != nil {
//		return fmt.Errorf("whatsapp marshal template request: %w", err)
//	}
//
//	req, err := http.NewRequestWithContext(
//		ctx,
//		http.MethodPost,
//		s.client.messagesURL(),
//		bytes.NewReader(bodyBytes),
//	)
//	if err != nil {
//		return fmt.Errorf("whatsapp build request: %w", err)
//	}
//
//	resp, err := s.client.doRequest(req)
//	if err != nil {
//		return fmt.Errorf("whatsapp send request: %w", err)
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode >= 400 {
//		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 4096))
//		if readErr != nil {
//			return fmt.Errorf("whatsapp api status %d: read body: %w", resp.StatusCode, readErr)
//		}
//
//		return fmt.Errorf("whatsapp api status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
//	}
//	return nil
//}
