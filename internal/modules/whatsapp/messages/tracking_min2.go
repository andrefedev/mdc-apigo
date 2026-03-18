package messages

import "time"

type MessageCategoryMin2 string

const (
	MessageCategoryAuthMin2      MessageCategoryMin2 = "auth"
	MessageCategoryUtilityMin2   MessageCategoryMin2 = "utility"
	MessageCategoryMarketingMin2 MessageCategoryMin2 = "marketing"
)

type MessageStatusMin2 string

const (
	MessageStatusAcceptedMin2  MessageStatusMin2 = "accepted"
	MessageStatusSentMin2      MessageStatusMin2 = "sent"
	MessageStatusDeliveredMin2 MessageStatusMin2 = "delivered"
	MessageStatusReadMin2      MessageStatusMin2 = "read"
	MessageStatusFailedMin2    MessageStatusMin2 = "failed"
)

type MessageEventTypeMin2 string

const (
	MessageEventTypeAcceptedMin2 MessageEventTypeMin2 = "accepted"
	MessageEventTypeStatusMin2   MessageEventTypeMin2 = "status"
	MessageEventTypeFailedMin2   MessageEventTypeMin2 = "failed"
)

type TrackedMessageMin2 struct {
	Ref          string
	UserRef      *string
	WaID         *string
	ToPhone      string
	MessageID    string
	Category     MessageCategoryMin2
	TemplateName *string
	Status       MessageStatusMin2
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TrackedMessageEventMin2 struct {
	Ref       string
	MessageID string
	EventType MessageEventTypeMin2
	Status    MessageStatusMin2
	CreatedAt time.Time
}
