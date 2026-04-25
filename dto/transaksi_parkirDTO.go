package dto

import (
	"time"
)

type TransaksiParkirRequestDTO struct {
	QRPayload      string `json:"qr_payload" binding:"required"`
	JenisKendaraan string `json:"jenis_kendaraan" binding:"required"`
}

type TransaksiParkirResponseDTO struct {
	ID             string                   `json:"id"`
	KodeTransaksi  string                   `json:"kode_transaksi"`
	Lokasi         *LokasiParkirResponseDTO `json:"lokasi"`
	Petugas        *UserResponseDTO         `json:"petugas"`
	JenisKendaraan string                   `json:"jenis_kendaraan"`
	Tarif          int64                    `json:"tarif"`
	SaldoSebelum   int64                    `json:"saldo_sebelum"`
	SaldoSesudah   int64                    `json:"saldo_sesudah"`
	Status         string                   `json:"status"`
	Waktu          time.Time                `json:"waktu"`
}
