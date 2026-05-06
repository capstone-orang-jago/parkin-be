-- +migrate Up

CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "citext";

CREATE TYPE user_role AS ENUM ('pengguna', 'petugas', 'dishub', 'superadmin');
CREATE TYPE vehicle_type AS ENUM ('motor', 'mobil');
CREATE TYPE transaction_status AS ENUM ('sukses', 'gagal', 'pending');
CREATE TYPE topup_status AS ENUM ('pending', 'sukses', 'gagal', 'expired');

CREATE TABLE users (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name    VARCHAR(100) NOT NULL,
    telephone     VARCHAR(20)  NOT NULL UNIQUE,
    email          VARCHAR(100),
    password  VARCHAR(255) NOT NULL,
    role           user_role    NOT NULL DEFAULT 'pengguna',
    saldo          BIGINT       NOT NULL DEFAULT 0 CHECK (saldo >= 0),
    is_verified    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_users_role ON users(role);
CREATE UNIQUE INDEX idx_users_email 
ON users(LOWER(email)) 
WHERE email IS NOT NULL;

CREATE TABLE otp_codes (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code    VARCHAR(255) NOT NULL,
    expired_at   TIMESTAMPTZ  NOT NULL,
    attempts     SMALLINT     NOT NULL DEFAULT 0,
    is_used      BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_otp_user_active 
ON otp_codes(user_id, is_used, expired_at);

CREATE TABLE lokasi_parkir (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama_lokasi VARCHAR(100) NOT NULL,
    alamat      TEXT,
    petugas_id  UUID REFERENCES users(id) ON DELETE SET NULL,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_lokasi_petugas 
ON lokasi_parkir(petugas_id) 
WHERE is_active;

CREATE TABLE qr_sessions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lokasi_id   UUID NOT NULL REFERENCES lokasi_parkir(id),
    petugas_id  UUID NOT NULL REFERENCES users(id),
    payload     TEXT NOT NULL,
    hmac   VARCHAR(64) NOT NULL,
    expired_at  TIMESTAMPTZ NOT NULL,
    is_revoked  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_qr_petugas_active 
ON qr_sessions(petugas_id, is_revoked, expired_at);

CREATE INDEX idx_qr_payload 
ON qr_sessions USING hash (payload);

CREATE TABLE transaksi_parkir (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kode_transaksi  VARCHAR(30) NOT NULL UNIQUE,
    pengguna_id     UUID NOT NULL REFERENCES users(id),
    petugas_id      UUID NOT NULL REFERENCES users(id),
    lokasi_id       UUID NOT NULL REFERENCES lokasi_parkir(id),
    qr_session_id   UUID REFERENCES qr_sessions(id),
    jenis_kendaraan vehicle_type NOT NULL,
    tarif           BIGINT NOT NULL CHECK (tarif > 0),
    saldo_sebelum   BIGINT NOT NULL,
    saldo_sesudah   BIGINT NOT NULL,
    status          transaction_status NOT NULL DEFAULT 'sukses',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_trx_pengguna_date 
ON transaksi_parkir(pengguna_id, created_at DESC);

CREATE INDEX idx_trx_petugas_date  
ON transaksi_parkir(petugas_id, created_at DESC);

CREATE INDEX idx_trx_lokasi_date   
ON transaksi_parkir(lokasi_id, created_at DESC);

CREATE TABLE topup_saldo (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pengguna_id             UUID NOT NULL REFERENCES users(id),
    nominal                 BIGINT NOT NULL CHECK (nominal >= 10000),
    metode_pembayaran       VARCHAR(50) NOT NULL,
    midtrans_order_id       VARCHAR(80) NOT NULL UNIQUE,
    midtrans_transaction_id VARCHAR(100) UNIQUE,
    status                  topup_status NOT NULL DEFAULT 'pending',
    snap_token              TEXT,
    snap_redirect_url       TEXT,
    raw_callback            JSONB,
    paid_at                 TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_topup_user_status 
ON topup_saldo(pengguna_id, status, created_at DESC);

-- +migrate Down
DROP TABLE IF EXISTS topup_saldo;
DROP TABLE IF EXISTS transaksi_parkir;
DROP TABLE IF EXISTS qr_sessions;
DROP TABLE IF EXISTS lokasi_parkir;
DROP TABLE IF EXISTS otp_codes;
DROP TABLE IF EXISTS users;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        DROP TYPE user_role;
    END IF;
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'vehicle_type') THEN
        DROP TYPE vehicle_type;
    END IF;
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'transaction_status') THEN
        DROP TYPE transaction_status;
    END IF;
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'topup_status') THEN
        DROP TYPE topup_status;
    END IF;
END $$;
