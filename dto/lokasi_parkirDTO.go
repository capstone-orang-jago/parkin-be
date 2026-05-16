package dto

import "time"

type LokasiParkirRequestDTO struct {
	NamaLokasi string `json:"nama_lokasi" binding:"required"`
	Alamat     string `json:"alamat" binding:"required"`
	PetugasID  string `json:"petugas_id" binding:"required"`
	IsActive   bool   `json:"is_active"`
}

type LokasiParkirResponseDTO struct {
	ID         string           `json:"id"`
	NamaLokasi string           `json:"nama_lokasi"`
	Alamat     string           `json:"alamat"`
	PetugasID  string           `json:"petugas_id"`
	IsActive   bool             `json:"is_active"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	Petugas    *UserResponseDTO `json:"petugas"`
}
