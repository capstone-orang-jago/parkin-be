package domain

import (
	"context"
	"time"

	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/gin-gonic/gin"
)

type TransaksiParkir struct {
	ID             string
	KodeTransaksi  string
	PenggunaID     string
	PetugasID      string
	LokasiID       string
	JenisKendaraan string
	Tarif          int64
	SaldoSebelum   int64
	SaldoSesudah   int64
	Status         string
	CreatedAt      time.Time
	Pengguna       *User
	LokasiParkir   *LokasiParkir
}

type TransaksiQuery struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Filter string `form:"filter" binding:"omitempty,oneof=semua minggu bulan"`
}

type (
	TransaksiParkirRepository interface {
		CreateTransaksiParkir(ctx context.Context, transaksi *TransaksiParkir) error
		GetTransaksiParkirByID(ctx context.Context, id string) (*TransaksiParkir, error)
		GetAllTransaksiParkirByUserID(ctx context.Context, userID, role string, page, limit int, filter string) ([]*TransaksiParkir, error)
	}

	TransaksiParkirUsecase interface {
		CreateTransaksiParkir(ctx context.Context, req *dto.TransaksiParkirRequestDTO, penggunaID string) (res *dto.TransaksiParkirResponseDTO, err error)
		GetTransaksiParkirByID(ctx context.Context, id string) (*dto.TransaksiParkirResponseDTO, error)
		GetAllTransaksiParkirByUserID(ctx context.Context, userID, role string, page, limit int, filter string) ([]*dto.TransaksiParkirResponseDTO, error)
	}

	TransaksiParkirHandler interface {
		CreateTransaksiParkir(c *gin.Context)
		GetTransaksiParkirByID(c *gin.Context)
		GetAllTransaksiParkirByUserID(c *gin.Context)
	}
)
