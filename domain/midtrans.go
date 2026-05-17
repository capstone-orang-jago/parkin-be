package domain

import (
	"context"

	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/gin-gonic/gin"
)

type (
	MidtransUsecase interface {
		InitiateTopUp(ctx context.Context, userID string, req dto.TopUpRequestDTO) (*dto.TopUpResponseDTO, error)
		TopUpSaldo(ctx context.Context, cb dto.MidtransCallbackDTO) (bool, error)
		GetTopUpStatus(ctx context.Context, midtransOrderID string) (*dto.TopUpResponseDTO, error)
	}

	MidtransHandler interface {
		CreateTopUp(c *gin.Context)
		MidtransCallback(c *gin.Context)
		GetTopUpStatus(c *gin.Context)
	}
)
