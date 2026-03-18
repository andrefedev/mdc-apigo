package messages

import "time"

type MessageCategory2 string

const (
	MessageCategoryAuth2      MessageCategory2 = "auth"
	MessageCategoryUtility2   MessageCategory2 = "utility"
	MessageCategoryMarketing2 MessageCategory2 = "marketing"
)

type MessageStatus2 string

const (
	MessageStatusQueued2    MessageStatus2 = "queued"
	MessageStatusAccepted2  MessageStatus2 = "accepted"
	MessageStatusSent2      MessageStatus2 = "sent"
	MessageStatusDelivered2 MessageStatus2 = "delivered"
	MessageStatusRead2      MessageStatus2 = "read"
	MessageStatusFailed2    MessageStatus2 = "failed"
)

type MessageEventType2 string

const (
	MessageEventTypeQueued2     MessageEventType2 = "queued"
	MessageEventTypeAPIAccept2  MessageEventType2 = "api_accept"
	MessageEventTypeSendError2  MessageEventType2 = "send_error"
	MessageEventTypeStatusSync2 MessageEventType2 = "status_sync"
)

type LedgerMessage2 struct {
	Ref                 string
	Category            MessageCategory2
	Source              string
	TemplateName        string
	TemplateLanguage    string
	ToPhone             string
	UserRef             *string
	ProviderMessageID   *string
	ProviderContactWaID *string
	Status              MessageStatus2
	ProviderStatusRaw   *string
	ProviderErrorCode   *string
	ProviderErrorMsg    *string
	RequestPayloadJSON  []byte
	ResponsePayloadJSON []byte
	ContextJSON         []byte
	QueuedAt            time.Time
	AcceptedAt          *time.Time
	SentAt              *time.Time
	DeliveredAt         *time.Time
	ReadAt              *time.Time
	FailedAt            *time.Time
	LastWebhookAt       *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type LedgerMessageEvent2 struct {
	Ref               string
	MessageRef        string
	ProviderMessageID *string
	EventType         MessageEventType2
	Status            *MessageStatus2
	ProviderStatusRaw *string
	ProviderErrorCode *string
	ProviderErrorMsg  *string
	PayloadJSON       []byte
	CreatedAt         time.Time
}
