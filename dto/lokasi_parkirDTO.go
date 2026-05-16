package dto

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
	IsActive   bool             `json:"is_active"`
	Petugas    *UserResponseDTO `json:"petugas"`
}
