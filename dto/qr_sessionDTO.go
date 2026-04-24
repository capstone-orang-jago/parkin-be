package dto

import "time"

type QRSessionResponseDTO struct {
	LocationID string    `json:"location_id"`
	PetugasID  string    `json:"petugas_id"`
	Payload    string    `json:"payload"`
	HMAC       string    `json:"hmac"`
	ExpiresAt  time.Time `json:"expires_at"`
	IsRevoked  bool      `json:"is_revoked"`
}
