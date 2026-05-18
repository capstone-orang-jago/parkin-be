package dto

import "time"

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

type AdminSummaryDTO struct {
	TotalTransactions   int64 `json:"total_transactions"`
	TotalRevenueRp      int64 `json:"total_revenue_rp"`
	MotorCount          int64 `json:"motor_count"`
	MobilCount          int64 `json:"mobil_count"`
	ActiveLocations     int64 `json:"active_locations"`
	ActiveOfficersToday int64 `json:"active_officers_today"`
}

type AdminTrendDTO struct {
	Date             string `json:"date"`  // Format "YYYY-MM-DD"
	Label            string `json:"label"` // Format hari: "Sen", "Sel", dst.
	TransactionCount int64  `json:"transaction_count"`
	RevenueRp        int64  `json:"revenue_rp"`
}

type AdminRecentTransactionDTO struct {
	ID           string    `json:"id"`
	Time         time.Time `json:"time"`
	LocationName string    `json:"location_name"`
	VehicleType  string    `json:"vehicle_type"`
	AmountRp     int64     `json:"amount_rp"`
	Status       string    `json:"status"`
}

type LocationPerformanceDTO struct {
	LocationName      string `json:"location_name"`
	TransactionsToday int64  `json:"transactions_today"`
	RevenueTodayRp    int64  `json:"revenue_today_rp"`
	OfficersOnDuty    string `json:"officers_on_duty"` // Nama-nama petugas digabung string
}

type AdminDashboardResponseDTO struct {
	Range               string                       `json:"range"`
	PeriodLabel         string                       `json:"period_label"`
	Summary             AdminSummaryDTO              `json:"summary"`
	Trend               []*AdminTrendDTO             `json:"trend"`
	RecentTransactions  []*AdminRecentTransactionDTO `json:"recent_transactions"`
	LocationPerformance []*LocationPerformanceDTO    `json:"location_performance"`
}
