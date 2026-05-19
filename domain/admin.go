package domain

import (
	"context"

	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/gin-gonic/gin"
)

// AdminTransaksiQuery menampung semua parameter filter dari URL query string
type AdminTransaksiQuery struct {
	DateFrom   string `form:"date_from"`
	DateTo     string `form:"date_to"`
	LocationID string `form:"location_id"`
	OfficerID  string `form:"officer_id"`
	Vehicle    string `form:"vehicle"`
	Status     string `form:"status"`
	Page       int    `form:"page" binding:"omitempty,min=1"`
	Limit      int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

type AdminReportQuery struct {
	DateFrom   string `form:"date_from"`
	DateTo     string `form:"date_to"`
	LocationID string `form:"location_id"`
	Vehicle    string `form:"vehicle"`
}

type (
	AdminRepository interface {
		GetAllTransactions(ctx context.Context, q AdminTransaksiQuery) ([]*TransaksiParkir, int64, error)
		GetTransactionByID(ctx context.Context, id string) (*TransaksiParkir, error)
		GetFilterOptions(ctx context.Context) ([]*dto.FilterOptionItem, []*dto.FilterOptionItem, error)
		GetAllLocations(ctx context.Context) ([]*LokasiParkir, error)
		GetLocationStatsByID(ctx context.Context, id string) (*dto.AdminLocationDetailDTO, error)
		GetAllOfficers(ctx context.Context) ([]*User, error)
		GetOfficerPerformanceByID(ctx context.Context, id string) (*dto.AdminOfficerDetailDTO, error)
		GetRevenueReport(ctx context.Context, q AdminReportQuery) (*dto.AdminReportResponseDTO, error)
	}

	AdminUsecase interface {
		GetAllTransactions(ctx context.Context, q AdminTransaksiQuery) ([]*dto.TransaksiParkirResponseDTO, int64, error)
		GetTransactionDetail(ctx context.Context, id string) (*dto.TransaksiParkirResponseDTO, error)
		GetFilterOptions(ctx context.Context) (*dto.AdminFilterOptionsResponseDTO, error)
		GetAllLocations(ctx context.Context) ([]*dto.LokasiParkirResponseDTO, error)
		GetLocationDetail(ctx context.Context, id string) (*dto.AdminLocationDetailDTO, error)
		GetAllOfficers(ctx context.Context) ([]*dto.UserResponseDTO, error)
		GetOfficerDetail(ctx context.Context, id string) (*dto.AdminOfficerDetailDTO, error)
		GetRevenueReport(ctx context.Context, q AdminReportQuery) (*dto.AdminReportResponseDTO, error)
	}

	AdminHandler interface {
		GetAllTransactions(c *gin.Context)
		GetTransactionDetail(c *gin.Context)
		GetFilterOptions(c *gin.Context)
		GetAllLocations(c *gin.Context)
		GetLocationDetail(c *gin.Context)
		GetAllOfficers(c *gin.Context)
		GetOfficerDetail(c *gin.Context)
		GetRevenueReport(c *gin.Context)
	}
)
