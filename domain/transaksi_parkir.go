package domain

import (
	"context"

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
}

type (
	TransaksiParkirRepository interface {
		CreateTransaksiParkir(ctx context.Context, transaksi *TransaksiParkir) error
		GetTransaksiParkirByID(ctx context.Context, id string) (*TransaksiParkir, error)
		GetAllTransaksiParkirByPenggunaID(ctx context.Context, penggunaID string) ([]*TransaksiParkir, error)
	}

	TransaksiParkirUsecase interface {
		CreateTransaksiParkir(ctx context.Context, req *dto.TransaksiParkirRequestDTO, penggunaID string) (res *dto.TransaksiParkirResponseDTO, err error)
		GetTransaksiParkirByID(ctx context.Context, id string) (*dto.TransaksiParkirResponseDTO, error)
		GetAllTransaksiParkirByPenggunaID(ctx context.Context, penggunaID string) ([]*dto.TransaksiParkirResponseDTO, error)
	}

	TransaksiParkirHandler interface {
		CreateTransaksiParkir(c *gin.Context)
		GetTransaksiParkirByID(c *gin.Context)
		GetAllTransaksiParkirByPenggunaID(c *gin.Context)
	}
)
