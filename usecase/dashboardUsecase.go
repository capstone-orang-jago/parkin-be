package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/sirupsen/logrus"
)

type dashboardUsecase struct {
	dashboardRepo domain.DashboardRepository
	transaksiRepo domain.TransaksiParkirRepository
	userRepo      domain.UserRepository
}

func NewDashboardUsecase(dashboardRepo domain.DashboardRepository, transaksiRepo domain.TransaksiParkirRepository, userRepo domain.UserRepository) domain.DashboardUsecase {
	return &dashboardUsecase{
		dashboardRepo: dashboardRepo,
		transaksiRepo: transaksiRepo,
		userRepo:      userRepo,
	}
}

func (usecase *dashboardUsecase) GetOfficerDashboard(ctx context.Context, officerID string) (*dto.OfficerDashboardResponseDTO, error) {
	userProfile, err := usecase.userRepo.GetUserProfile(ctx, officerID, "petugas")
	if err != nil || userProfile == nil || userProfile.LokasiParkir == nil {
		return nil, errors.New("petugas tidak ditemukan atau tidak memiliki lokasi parkir")
	}

	lokasiID := userProfile.LokasiParkir.ID

	totalMotor, totalMobil, totalPendapatan, err := usecase.dashboardRepo.GetWeeklyDashboardData(ctx, lokasiID)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetOfficerDashboard: failed to get weekly stats")
		return nil, err
	}

	recentTx, err := usecase.transaksiRepo.GetAllTransaksiParkirByUserID(ctx, officerID, "petugas", 1, 5, "semua")
	if err != nil {
		logrus.WithError(err).Error("usecase.GetOfficerDashboard: failed to get recent transactions")
		return nil, err
	}

	var txDTOs []*dto.OfficerRecentTransactionDTO
	for _, t := range recentTx {
		txDTOs = append(txDTOs, &dto.OfficerRecentTransactionDTO{
			ID:             t.ID,
			Waktu:          helpers.FormatTime(t.CreatedAt),
			JenisKendaraan: t.JenisKendaraan,
			Tarif:          t.Tarif,
		})
	}

	return &dto.OfficerDashboardResponseDTO{
		NamaPetugas:     userProfile.Name,
		NamaLokasi:      userProfile.LokasiParkir.NamaLokasi,
		Alamat:          userProfile.LokasiParkir.Alamat,
		TotalMotor:      totalMotor,
		TotalMobil:      totalMobil,
		TotalPendapatan: totalPendapatan,
		Transactions:    txDTOs,
	}, nil
}

func (usecase *dashboardUsecase) GetAdminDashboard(ctx context.Context, filterRange string) (*dto.AdminDashboardResponseDTO, error) {
	now := helpers.GetLocalTime()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var kpiStartTime time.Time
	var trendDays int
	var periodLabel string

	indoMonths := []string{"", "Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"}

	switch filterRange {
	case "week":
		kpiStartTime = startOfToday.AddDate(0, 0, -6)
		trendDays = 7
		periodLabel = fmt.Sprintf("7 Hari Terakhir · %s - %s", kpiStartTime.Format("02 Jan 2006"), now.Format("02 Jan 2006"))
	case "month":
		kpiStartTime = startOfToday.AddDate(0, 0, -29)
		trendDays = 30
		periodLabel = fmt.Sprintf("30 Hari Terakhir · %s - %s", kpiStartTime.Format("02 Jan 2006"), now.Format("02 Jan 2006"))
	default: // "today"
		filterRange = "today"
		kpiStartTime = startOfToday
		trendDays = 7
		periodLabel = fmt.Sprintf("Hari ini · %d %s %d", now.Day(), indoMonths[now.Month()], now.Year())
	}

	trendStartTime := startOfToday.AddDate(0, 0, -(trendDays - 1))

	dbStartTime := kpiStartTime
	if trendStartTime.Before(kpiStartTime) {
		dbStartTime = trendStartTime
	}

	rawTransactions, err := usecase.dashboardRepo.GetRawTransactionsForAdmin(ctx, dbStartTime)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetAdminDashboard: failed to get raw transactions")
		return nil, err
	}

	var summary dto.AdminSummaryDTO
	trendMap := make(map[string]*dto.AdminTrendDTO)
	perfMap := make(map[string]*dto.LocationPerformanceDTO)
	officerDutyMap := make(map[string]map[string]bool)
	var recentTransactions []*dto.AdminRecentTransactionDTO

	indoDays := map[time.Weekday]string{
		time.Sunday: "Min", time.Monday: "Sen", time.Tuesday: "Sel",
		time.Wednesday: "Rab", time.Thursday: "Kam", time.Friday: "Jum", time.Saturday: "Sab",
	}

	for _, tx := range rawTransactions {
		localTxTime := tx.CreatedAt.In(now.Location())
		txDateStr := localTxTime.Format("2006-01-02")

		if !localTxTime.Before(trendStartTime) {
			if _, exists := trendMap[txDateStr]; !exists {
				trendMap[txDateStr] = &dto.AdminTrendDTO{
					Date:             txDateStr,
					TransactionCount: 0,
					RevenueRp:        0,
				}
			}
			trendMap[txDateStr].TransactionCount++
			trendMap[txDateStr].RevenueRp += tx.Tarif
		}
		if !localTxTime.Before(kpiStartTime) {
			summary.TotalTransactions++
			summary.TotalRevenueRp += tx.Tarif

			switch tx.JenisKendaraan {
			case "motor":
				summary.MotorCount++
			case "mobil":
				summary.MobilCount++
			}

			if tx.LokasiParkir != nil && tx.LokasiParkir.NamaLokasi != "" {
				locName := tx.LokasiParkir.NamaLokasi

				if _, exists := perfMap[locName]; !exists {
					perfMap[locName] = &dto.LocationPerformanceDTO{
						LocationName:      locName,
						TransactionsToday: 0,
						RevenueTodayRp:    0,
					}
					officerDutyMap[locName] = make(map[string]bool)
				}
				perfMap[locName].TransactionsToday++
				perfMap[locName].RevenueTodayRp += tx.Tarif

				if tx.LokasiParkir.Petugas != nil && tx.LokasiParkir.Petugas.Name != "" {
					officerDutyMap[locName][tx.LokasiParkir.Petugas.Name] = true
				}
			}
		}
	}

	for i := 0; i < len(rawTransactions) && i < 10; i++ {
		tx := rawTransactions[i]
		locName := "-"
		if tx.LokasiParkir != nil && tx.LokasiParkir.NamaLokasi != "" {
			locName = tx.LokasiParkir.NamaLokasi
		}
		recentTransactions = append(recentTransactions, &dto.AdminRecentTransactionDTO{
			ID:           tx.ID,
			Time:         tx.CreatedAt.In(now.Location()),
			LocationName: locName,
			VehicleType:  tx.JenisKendaraan,
			AmountRp:     tx.Tarif,
			Status:       tx.Status,
		})
	}

	var trendSlice []*dto.AdminTrendDTO
	for d := 0; d < trendDays; d++ {
		dayDate := trendStartTime.AddDate(0, 0, d)
		dayStr := dayDate.Format("2006-01-02")

		label := indoDays[dayDate.Weekday()]
		if trendDays > 7 {
			label = dayDate.Format("02 Jan")
		}

		if val, exists := trendMap[dayStr]; exists {
			val.Label = label
			trendSlice = append(trendSlice, val)
		} else {
			trendSlice = append(trendSlice, &dto.AdminTrendDTO{
				Date:             dayStr,
				Label:            label,
				TransactionCount: 0,
				RevenueRp:        0,
			})
		}
	}

	var performanceSlice []*dto.LocationPerformanceDTO
	uniqueOfficersGlobally := make(map[string]bool)

	for locName, perfData := range perfMap {
		var officers []string
		for officerName := range officerDutyMap[locName] {
			officers = append(officers, officerName)
			uniqueOfficersGlobally[officerName] = true
		}

		perfData.OfficersOnDuty = "Tidak ada petugas"
		if len(officers) > 0 {
			perfData.OfficersOnDuty = strings.Join(officers, ", ")
		}
		performanceSlice = append(performanceSlice, perfData)
	}

	summary.ActiveLocations = int64(len(perfMap))
	summary.ActiveOfficersToday = int64(len(uniqueOfficersGlobally))

	return &dto.AdminDashboardResponseDTO{
		Range:               filterRange,
		PeriodLabel:         periodLabel,
		Summary:             summary,
		Trend:               trendSlice,
		RecentTransactions:  recentTransactions,
		LocationPerformance: performanceSlice,
	}, nil
}
