-- +migrate Up
ALTER TABLE qr_sessions ADD COLUMN IF NOT EXISTS last_used_at TIMESTAMPTZ;

-- +migrate Down
ALTER TABLE qr_sessions DROP COLUMN IF EXISTS last_used_at;
