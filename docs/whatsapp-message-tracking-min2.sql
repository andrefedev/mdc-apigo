CREATE TABLE IF NOT EXISTS whatsapp_messages_min (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_ref uuid NULL,
    wa_id varchar(64) NULL,
    to_phone varchar(32) NOT NULL,
    message_id varchar(255) NOT NULL UNIQUE,
    category varchar(32) NOT NULL,
    template_name varchar(128) NULL,
    status varchar(32) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_whatsapp_messages_min_category_created_at
    ON whatsapp_messages_min(category, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_whatsapp_messages_min_status_created_at
    ON whatsapp_messages_min(status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_whatsapp_messages_min_wa_id_created_at
    ON whatsapp_messages_min(wa_id, created_at DESC);

CREATE TABLE IF NOT EXISTS whatsapp_message_events_min (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    message_id varchar(255) NOT NULL,
    event_type varchar(32) NOT NULL,
    status varchar(32) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_whatsapp_message_events_min_message_id_created_at
    ON whatsapp_message_events_min(message_id, created_at DESC);
