# syntax=docker/dockerfile:1

# ── Stage 1: build ──────────────────────────────────────────────
# Hanya dipakai untuk compile. Go toolchain + module cache + source
# (±350 MB) semuanya dibuang setelah binary jadi.
FROM golang:1.25-alpine AS build

WORKDIR /src

# Download dependency dulu, terpisah dari source, supaya layer ini
# ter-cache dan build ulang tidak download ulang.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 → binary statis, tidak butuh libc di runtime.
# -ldflags "-s -w" → buang tabel simbol/debug, binary lebih kecil.
# -trimpath → path build tidak ikut ke dalam binary.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/parkin-be .

# ── Stage 2: runtime ────────────────────────────────────────────
# Hanya yang dibutuhkan untuk jalan: binary + file migrasi/seeder +
# sertifikat HTTPS + data zona waktu. Tidak ada Go, tidak ada source.
FROM alpine:3.21

# ca-certificates: backend memanggil Midtrans & WhatsApp lewat HTTPS.
# tzdata: timestamp log & transaksi sesuai zona waktu.
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=build /out/parkin-be ./parkin-be

# Wajib ikut: migrate:up dan seeder:all membaca file SQL dari sini
# (database/migrations dan database/seeders) saat runtime.
COPY --from=build /src/database ./database

# Alpine sengaja dipilih, bukan scratch (0 MB):
# scratch tidak punya shell, jadi Terminal Coolify tidak bisa masuk
# ke container untuk menjalankan ./parkin-be migrate:up.
# Alpine hanya menambah ±7 MB tapi shell tetap ada.
USER nobody

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=60s --retries=3 \
  CMD wget -q --spider http://localhost:8080/api/v1/health || exit 1

CMD ["./parkin-be", "serve"]
