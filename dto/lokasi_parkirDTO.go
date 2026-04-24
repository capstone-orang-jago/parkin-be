package dto

import "time"

type LokasiParkirRequestDTO struct {
	NamaLokasi string
	Alamat     string
	PetugasID  string
	IsActive   bool
}

type LokasiParkirResponseDTO struct {
	ID         string
	NamaLokasi string
	Alamat     string
	PetugasID  string
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
