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

func (usecase *dashboardUsecase) GetAdminDashboard(ctx context.Context) (*dto.AdminDashboardResponseDTO, error) {
	now := helpers.GetLocalTime()
	todayStr := now.Format("2006-01-02")

	sevenDaysAgo := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -6)

	rawTransactions, err := usecase.dashboardRepo.GetRawTransactionsForAdmin(ctx, sevenDaysAgo)
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

	for i, tx := range rawTransactions {

		localTxTime := tx.CreatedAt.In(now.Location())
		txDateStr := localTxTime.Format("2006-01-02")

		if _, exists := trendMap[txDateStr]; !exists {
			trendMap[txDateStr] = &dto.AdminTrendDTO{
				Date:             txDateStr,
				Label:            indoDays[localTxTime.Weekday()],
				TransactionCount: 0,
				RevenueRp:        0,
			}
		}
		trendMap[txDateStr].TransactionCount++
		trendMap[txDateStr].RevenueRp += tx.Tarif

		if txDateStr == todayStr {
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

		if i < 10 {
			locName := "-"
			if tx.LokasiParkir != nil && tx.LokasiParkir.NamaLokasi != "" {
				locName = tx.LokasiParkir.NamaLokasi
			}
			recentTransactions = append(recentTransactions, &dto.AdminRecentTransactionDTO{
				ID:           tx.ID,
				Time:         localTxTime,
				LocationName: locName,
				VehicleType:  tx.JenisKendaraan,
				AmountRp:     tx.Tarif,
				Status:       tx.Status,
			})
		}
	}

	var trendSlice []*dto.AdminTrendDTO
	for d := 0; d < 7; d++ {
		dayDate := sevenDaysAgo.AddDate(0, 0, d)
		dayStr := dayDate.Format("2006-01-02")

		if val, exists := trendMap[dayStr]; exists {
			trendSlice = append(trendSlice, val)
		} else {
			trendSlice = append(trendSlice, &dto.AdminTrendDTO{
				Date:             dayStr,
				Label:            indoDays[dayDate.Weekday()],
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

	indoMonths := []string{"", "Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"}
	periodLabel := fmt.Sprintf("Hari ini · %d %s %d", now.Day(), indoMonths[now.Month()], now.Year())

	return &dto.AdminDashboardResponseDTO{
		Range:               "today",
		PeriodLabel:         periodLabel,
		Summary:             summary,
		Trend:               trendSlice,
		RecentTransactions:  recentTransactions,
		LocationPerformance: performanceSlice,
	}, nil
}
