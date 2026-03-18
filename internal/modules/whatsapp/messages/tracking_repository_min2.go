package messages

import "context"

type CreateTrackedMessageMin2 struct {
	UserRef      *string
	WaID         *string
	ToPhone      string
	MessageID    string
	Category     MessageCategoryMin2
	TemplateName *string
	Status       MessageStatusMin2
}

type UpdateTrackedMessageStatusMin2 struct {
	MessageID string
	Status    MessageStatusMin2
}

type CreateTrackedMessageEventMin2 struct {
	MessageID string
	EventType MessageEventTypeMin2
	Status    MessageStatusMin2
}

type TrackedMessageRepositoryMin2 interface {
	CreateMessage(ctx context.Context, input *CreateTrackedMessageMin2) (*TrackedMessageMin2, error)
	UpdateStatus(ctx context.Context, input *UpdateTrackedMessageStatusMin2) error
	CreateEvent(ctx context.Context, input *CreateTrackedMessageEventMin2) error
	FindByMessageID(ctx context.Context, messageID string) (*TrackedMessageMin2, error)
	ListMessages(ctx context.Context) ([]*TrackedMessageMin2, error)
}
