CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS auth_sessions (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    uid uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash text NOT NULL UNIQUE,
    last_used_at timestamptz NULL,
    date_expires timestamptz NOT NULL,
    date_created timestamptz NOT NULL DEFAULT NOW(),
    date_revoked timestamptz NULL
);

CREATE INDEX IF NOT EXISTS idx_auth_sessions_uid ON auth_sessions (uid);
CREATE INDEX IF NOT EXISTS idx_auth_sessions_active_token ON auth_sessions (token_hash)
    WHERE date_revoked IS NULL;
