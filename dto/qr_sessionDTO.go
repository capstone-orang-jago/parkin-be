package dto

type QRSessionResponseDTO struct {
	ID           string                   `json:"id"`
	Payload      string                   `json:"payload"`
	HMAC         string                   `json:"hmac"`
	ExpiredAt    string                   `json:"expired_at"`
	IsRevoked    bool                     `json:"is_revoked"`
	LokasiParkir *LokasiParkirResponseDTO `json:"lokasi_parkir"`
}
