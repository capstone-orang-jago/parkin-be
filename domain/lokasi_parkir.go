package domain

import (
	"context"
	"time"

	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/gin-gonic/gin"
)

type (
	LokasiParkir struct {
		ID         string
		NamaLokasi string
		Alamat     string
		PetugasID  string
		IsActive   bool
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}

	LokasiParkirRepository interface {
		CreateLokasiParkir(ctx context.Context, lokasiParkir *LokasiParkir) error
		GetAllLokasiParkir(ctx context.Context) ([]*LokasiParkir, error)
		GetLokasiParkirByID(ctx context.Context, id string) (*LokasiParkir, error)
		GetLokasiParkirByPetugasID(ctx context.Context, petugasID string) (*LokasiParkir, error)
		UpdateLokasiParkir(ctx context.Context, lokasiParkir *LokasiParkir, id string) error
		DeleteLokasiParkir(ctx context.Context, id string) error
	}

	LokasiParkirUsecase interface {
		CreateLokasiParkir(ctx context.Context, req *dto.LokasiParkirRequestDTO) (res *dto.LokasiParkirResponseDTO, err error)
		GetAllLokasiParkir(ctx context.Context) ([]*dto.LokasiParkirResponseDTO, error)
		GetLokasiParkirByID(ctx context.Context, id string) (*dto.LokasiParkirResponseDTO, error)
		GetLokasiParkirByPetugasID(ctx context.Context, petugasID string) (*dto.LokasiParkirResponseDTO, error)
		UpdateLokasiParkir(ctx context.Context, req dto.LokasiParkirRequestDTO, id string) (res *dto.LokasiParkirResponseDTO, err error)
		DeleteLokasiParkir(ctx context.Context, id string) error
	}

	LokasiParkirHandler interface {
		CreateLokasiParkir(c *gin.Context)
		GetAllLokasiParkir(c *gin.Context)
		GetLokasiParkirByID(c *gin.Context)
		GetLokasiParkirByPetugasID(c *gin.Context)
		UpdateLokasiParkir(c *gin.Context)
		DeleteLokasiParkir(c *gin.Context)
	}
)
