CREATE TABLE IF NOT EXISTS whatsapp_messages (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    category varchar(32) NOT NULL,
    source varchar(64) NOT NULL,
    template_name varchar(128) NOT NULL,
    template_language varchar(32) NOT NULL DEFAULT 'es_CO',
    to_phone varchar(32) NOT NULL,
    user_ref uuid NULL,
    provider_message_id varchar(255) NULL,
    provider_contact_wa_id varchar(64) NULL,
    status varchar(32) NOT NULL,
    provider_status_raw varchar(64) NULL,
    provider_error_code varchar(64) NULL,
    provider_error_message text NULL,
    request_payload_json jsonb NULL,
    response_payload_json jsonb NULL,
    context_json jsonb NULL,
    queued_at timestamptz NOT NULL DEFAULT now(),
    accepted_at timestamptz NULL,
    sent_at timestamptz NULL,
    delivered_at timestamptz NULL,
    read_at timestamptz NULL,
    failed_at timestamptz NULL,
    last_webhook_at timestamptz NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_whatsapp_messages_provider_message_id
    ON waba_messages(provider_message_id)
    WHERE provider_message_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_whatsapp_messages_category_created_at
    ON waba_messages(category, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_whatsapp_messages_status_created_at
    ON waba_messages(status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_whatsapp_messages_to_phone_created_at
    ON waba_messages(to_phone, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_whatsapp_messages_user_ref_created_at
    ON waba_messages(user_ref, created_at DESC);

CREATE TABLE IF NOT EXISTS whatsapp_message_events (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    message_ref uuid NOT NULL REFERENCES waba_messages(id) ON DELETE CASCADE,
    provider_message_id varchar(255) NULL,
    event_type varchar(32) NOT NULL,
    status varchar(32) NULL,
    provider_status_raw varchar(64) NULL,
    provider_error_code varchar(64) NULL,
    provider_error_message text NULL,
    payload_json jsonb NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_whatsapp_message_events_message_ref_created_at
    ON whatsapp_message_events(message_ref, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_whatsapp_message_events_provider_message_id
    ON whatsapp_message_events(provider_message_id);
