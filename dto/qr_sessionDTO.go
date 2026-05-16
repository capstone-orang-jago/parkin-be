package dto

import "time"

type QRSessionResponseDTO struct {
	ID           string                   `json:"id"`
	LokasiID     string                   `json:"lokasi_id"`
	PetugasID    string                   `json:"petugas_id"`
	Payload      string                   `json:"payload"`
	HMAC         string                   `json:"hmac"`
	ExpiredAt    time.Time                `json:"expired_at"`
	IsRevoked    bool                     `json:"is_revoked"`
	LokasiParkir *LokasiParkirResponseDTO `json:"lokasi_parkir"`
	Petugas      *UserResponseDTO         `json:"petugas"`
}
