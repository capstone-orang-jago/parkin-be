package dto

type FilterOptionItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AdminFilterOptionsResponseDTO struct {
	Locations []*FilterOptionItem `json:"locations"`
	Officers  []*FilterOptionItem `json:"officers"`
	Vehicles  []string            `json:"vehicles"`
	Statuses  []string            `json:"statuses"`
}

type AdminLocationDetailDTO struct {
	ID              string                        `json:"id"`
	NamaLokasi      string                        `json:"nama_lokasi"`
	Alamat          string                        `json:"alamat"`
	IsActive        bool                          `json:"is_active"`
	TotalTransaksi  int64                         `json:"total_transaksi"`
	TotalPendapatan int64                         `json:"total_pendapatan"`
	RecentHistory   []*TransaksiParkirResponseDTO `json:"recent_history"`
}

type AdminOfficerDetailDTO struct {
	ID             string                        `json:"id"`
	Name           string                        `json:"name"`
	Email          string                        `json:"email"`
	Telephone      string                        `json:"telephone"`
	TotalDutyDays  int64                         `json:"total_duty_days"`
	TotalTxHandled int64                         `json:"total_tx_handled"`
	HistoryHandled []*TransaksiParkirResponseDTO `json:"history_handled"`
}

type ReportLocationBreakdown struct {
	LocationName string `json:"location_name"`
	TotalTx      int64  `json:"total_tx"`
	Revenue      int64  `json:"revenue"`
}

type AdminReportResponseDTO struct {
	TotalRevenue    int64                      `json:"total_revenue"`
	TotalTxGlobal   int64                      `json:"total_tx_global"`
	AveragePerTx    float64                    `json:"average_per_tx"`
	BreakdownBySpot []*ReportLocationBreakdown `json:"breakdown_by_spot"`
}

type PaginatedTransactionResponse struct {
	Total int64                         `json:"total"`
	Page  int                           `json:"page"`
	Limit int                           `json:"limit"`
	Items []*TransaksiParkirResponseDTO `json:"items"`
}
