package dto

type OfficerRecentTransactionDTO struct {
	ID             string `json:"id"`
	Waktu          string `json:"waktu"`
	JenisKendaraan string `json:"jenis_kendaraan"`
	Tarif          int64  `json:"tarif"`
}

type OfficerDashboardResponseDTO struct {
	NamaPetugas     string                         `json:"nama_petugas"`
	NamaLokasi      string                         `json:"nama_lokasi"`
	Alamat          string                         `json:"alamat"`
	TotalMotor      int                            `json:"total_motor"`
	TotalMobil      int                            `json:"total_mobil"`
	TotalPendapatan int64                          `json:"total_pendapatan"`
	Transactions    []*OfficerRecentTransactionDTO `json:"transactions"`
}
