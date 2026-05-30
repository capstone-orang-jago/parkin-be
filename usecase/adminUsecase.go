package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/sirupsen/logrus"
)

type adminUsecase struct {
	adminRepo domain.AdminRepository
}

func NewAdminUsecase(adminRepo domain.AdminRepository) domain.AdminUsecase {
	return &adminUsecase{adminRepo: adminRepo}
}

func (u *adminUsecase) GetAllTransactions(ctx context.Context, q domain.AdminTransaksiQuery) ([]*dto.TransaksiParkirResponseDTO, int64, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 {
		q.Limit = 10
	}

	txs, total, err := u.adminRepo.GetAllTransactions(ctx, q)
	if err != nil {
		logrus.WithError(err).Error("adminUsecase.GetAllTransactions: failed to get transactions from repository")
		return nil, 0, err
	}

	var dtos []*dto.TransaksiParkirResponseDTO
	for _, tx := range txs {
		res := &dto.TransaksiParkirResponseDTO{
			ID:             tx.ID,
			KodeTransaksi:  tx.KodeTransaksi,
			JenisKendaraan: tx.JenisKendaraan,
			Tarif:          tx.Tarif,
			SaldoSebelum:   tx.SaldoSebelum,
			SaldoSesudah:   tx.SaldoSesudah,
			Status:         tx.Status,
			CreatedAt:      helpers.FormatTime(tx.CreatedAt),
		}
		if tx.Pengguna != nil {
			res.Pengguna = &dto.UserResponseDTO{ID: tx.Pengguna.ID, Name: tx.Pengguna.Name, Email: tx.Pengguna.Email}
		}
		if tx.Petugas != nil {
			res.Petugas = &dto.UserResponseDTO{ID: tx.Petugas.ID, Name: tx.Petugas.Name, Email: tx.Petugas.Email}
		}
		if tx.LokasiParkir != nil {
			res.Lokasi = &dto.LokasiParkirResponseDTO{ID: tx.LokasiParkir.ID, NamaLokasi: tx.LokasiParkir.NamaLokasi}
		}
		dtos = append(dtos, res)
	}
	return dtos, total, nil
}

func (u *adminUsecase) GetTransactionDetail(ctx context.Context, id string) (*dto.TransaksiParkirResponseDTO, error) {
	tx, err := u.adminRepo.GetTransactionByID(ctx, id)
	if err != nil || tx == nil {
		logrus.WithError(err).Error("adminUsecase.GetTransactionDetail: failed to get transaction detail from repository")
		return nil, err
	}

	res := &dto.TransaksiParkirResponseDTO{
		ID:             tx.ID,
		KodeTransaksi:  tx.KodeTransaksi,
		JenisKendaraan: tx.JenisKendaraan,
		Tarif:          tx.Tarif,
		SaldoSebelum:   tx.SaldoSebelum,
		SaldoSesudah:   tx.SaldoSesudah,
		Status:         tx.Status,
		CreatedAt:      helpers.FormatTime(tx.CreatedAt),
	}
	if tx.Pengguna != nil {
		res.Pengguna = &dto.UserResponseDTO{ID: tx.Pengguna.ID, Name: tx.Pengguna.Name, Email: tx.Pengguna.Email}
	}
	if tx.Petugas != nil {
		res.Petugas = &dto.UserResponseDTO{ID: tx.Petugas.ID, Name: tx.Petugas.Name, Email: tx.Petugas.Email}
	}
	if tx.LokasiParkir != nil {
		res.Lokasi = &dto.LokasiParkirResponseDTO{ID: tx.LokasiParkir.ID, NamaLokasi: tx.LokasiParkir.NamaLokasi}
	}

	// if res.Pengguna != nil && len(res.Pengguna.ID) > 12 {
	// 	rawID := res.Pengguna.ID
	// 	res.Pengguna.ID = fmt.Sprintf("%s***%s", rawID[:5], rawID[len(rawID)-4:])
	// }
	return res, nil
}

func (u *adminUsecase) GetFilterOptions(ctx context.Context) (*dto.AdminFilterOptionsResponseDTO, error) {
	locs, offs, err := u.adminRepo.GetFilterOptions(ctx)
	if err != nil {
		logrus.WithError(err).Error("adminUsecase.GetFilterOptions: failed to get filter options from repository")
		return nil, err
	}

	return &dto.AdminFilterOptionsResponseDTO{
		Locations: locs,
		Officers:  offs,
		Vehicles:  []string{"motor", "mobil"},
		Statuses:  []string{"sukses", "gagal", "pending"},
	}, nil
}

func (u *adminUsecase) GetAllLocations(ctx context.Context) ([]*dto.LokasiParkirResponseDTO, error) {
	locations, err := u.adminRepo.GetAllLocations(ctx)
	if err != nil {
		logrus.WithError(err).Error("adminUsecase.GetAllLocations: failed to get locations from repository")
		return nil, err
	}

	var dtos []*dto.LokasiParkirResponseDTO
	for _, loc := range locations {
		locDTO := &dto.LokasiParkirResponseDTO{
			ID:         loc.ID,
			NamaLokasi: loc.NamaLokasi,
			Alamat:     loc.Alamat,
			IsActive:   loc.IsActive,
		}

		if loc.Petugas != nil {
			locDTO.Petugas = &dto.UserResponseDTO{
				ID:   loc.Petugas.ID,
				Name: loc.Petugas.Name,
			}
		}

		dtos = append(dtos, locDTO)
	}

	if dtos == nil {
		dtos = make([]*dto.LokasiParkirResponseDTO, 0)
	}
	return dtos, nil
}

func (u *adminUsecase) GetLocationDetail(ctx context.Context, id string) (*dto.AdminLocationDetailDTO, error) {
	stat, err := u.adminRepo.GetLocationStatsByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("adminUsecase.GetLocationDetail: failed to get location stats from repository")
		return nil, err
	}
	if stat == nil {
		return nil, errors.New("lokasi parkir tidak ditemukan")
	}

	recentTx, _, _ := u.GetAllTransactions(ctx, domain.AdminTransaksiQuery{LocationID: id, Limit: 5, Page: 1})

	stat.RecentHistory = recentTx
	if stat.RecentHistory == nil {
		stat.RecentHistory = make([]*dto.TransaksiParkirResponseDTO, 0)
	}

	return stat, nil
}

func (u *adminUsecase) GetAllOfficers(ctx context.Context) ([]*dto.UserResponseDTO, error) {
	officers, err := u.adminRepo.GetAllOfficers(ctx)
	if err != nil {
		logrus.WithError(err).Error("adminUsecase.GetAllOfficers: failed to get officers from repository")
		return nil, err
	}

	var dtos []*dto.UserResponseDTO
	for _, off := range officers {
		dtos = append(dtos, &dto.UserResponseDTO{
			ID:         off.ID,
			Name:       off.Name,
			Email:      off.Email,
			Telephone:  off.Telephone,
			Role:       off.Role,
			IsVerified: off.IsVerified,
		})
	}

	if dtos == nil {
		dtos = make([]*dto.UserResponseDTO, 0)
	}
	return dtos, nil
}

func (u *adminUsecase) GetOfficerDetail(ctx context.Context, id string) (*dto.AdminOfficerDetailDTO, error) {
	perf, err := u.adminRepo.GetOfficerPerformanceByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("adminUsecase.GetOfficerDetail: failed to get officer performance from repository")
		return nil, err
	}
	if perf == nil {
		return nil, errors.New("petugas tidak ditemukan")
	}

	recentTx, _, _ := u.GetAllTransactions(ctx, domain.AdminTransaksiQuery{OfficerID: id, Limit: 5, Page: 1})

	perf.HistoryHandled = recentTx
	if perf.HistoryHandled == nil {
		perf.HistoryHandled = make([]*dto.TransaksiParkirResponseDTO, 0)
	}

	return perf, nil
}

func (u *adminUsecase) GetRevenueReport(ctx context.Context, q domain.AdminReportQuery) (*dto.AdminReportResponseDTO, error) {
	if q.DateFrom == "" {
		q.DateFrom = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if q.DateTo == "" {
		q.DateTo = time.Now().Format("2006-01-02")
	}

	report, err := u.adminRepo.GetRevenueReport(ctx, q)
	if err != nil {
		logrus.WithError(err).Error("adminUsecase.GetRevenueReport: failed to get revenue report from repository")
		return nil, err
	}

	if report.BreakdownBySpot == nil {
		report.BreakdownBySpot = make([]*dto.ReportLocationBreakdown, 0)
	}

	return report, nil
}
