package messages

import "context"

type CreateLedgerMessageInput2 struct {
	Category           MessageCategory2
	Source             string
	TemplateName       string
	TemplateLanguage   string
	ToPhone            string
	UserRef            *string
	RequestPayloadJSON []byte
	ContextJSON        []byte
}

type MarkLedgerAcceptedInput2 struct {
	MessageRef          string
	ProviderMessageID   string
	ProviderContactWaID *string
	ProviderStatusRaw   *string
	ResponsePayloadJSON []byte
}

type MarkLedgerFailedInput2 struct {
	MessageRef          string
	ProviderErrorCode   *string
	ProviderErrorMsg    *string
	ResponsePayloadJSON []byte
}

type ApplyLedgerStatusInput2 struct {
	ProviderMessageID string
	Status            MessageStatus2
	ProviderStatusRaw *string
	ProviderErrorCode *string
	ProviderErrorMsg  *string
	PayloadJSON       []byte
}

type CreateLedgerEventInput2 struct {
	MessageRef        string
	ProviderMessageID *string
	EventType         MessageEventType2
	Status            *MessageStatus2
	ProviderStatusRaw *string
	ProviderErrorCode *string
	ProviderErrorMsg  *string
	PayloadJSON       []byte
}

type ListLedgerMessagesFilter2 struct {
	Category *MessageCategory2
	Status   *MessageStatus2
	ToPhone  *string
	UserRef  *string
}

type LedgerRepository2 interface {
	CreateMessage(ctx context.Context, input *CreateLedgerMessageInput2) (*LedgerMessage2, error)
	MarkAccepted(ctx context.Context, input *MarkLedgerAcceptedInput2) error
	MarkFailed(ctx context.Context, input *MarkLedgerFailedInput2) error
	ApplyStatus(ctx context.Context, input *ApplyLedgerStatusInput2) error
	CreateEvent(ctx context.Context, input *CreateLedgerEventInput2) error
	FindByProviderMessageID(ctx context.Context, providerMessageID string) (*LedgerMessage2, error)
	FindByRef(ctx context.Context, ref string) (*LedgerMessage2, error)
	ListMessages(ctx context.Context, filter *ListLedgerMessagesFilter2) ([]*LedgerMessage2, error)
}
