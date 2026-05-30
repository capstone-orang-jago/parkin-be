package domain

import (
	"context"
	"time"

	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/gin-gonic/gin"
)

type (
	AdminDashboardQuery struct {
		Range string `form:"range" binding:"omitempty,oneof=today week month"`
	}

	DashboardRepository interface {
		GetWeeklyDashboardData(ctx context.Context, lokasiID string) (motor, mobil int, pendapatan int64, err error)
		GetRawTransactionsForAdmin(ctx context.Context, startTime time.Time) ([]*TransaksiParkir, error)
	}

	DashboardUsecase interface {
		GetOfficerDashboard(ctx context.Context, officerID string) (res *dto.OfficerDashboardResponseDTO, err error)
		GetAdminDashboard(ctx context.Context, filterRange string) (*dto.AdminDashboardResponseDTO, error)
	}

	DashboardHandler interface {
		GetOfficerDashboard(c *gin.Context)
		GetAdminDashboard(c *gin.Context)
	}
)
