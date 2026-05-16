package domain

import (
	"context"
	"time"

	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/gin-gonic/gin"
)

type QRSession struct {
	ID           string
	LokasiID     string
	PetugasID    string
	Payload      string
	HMAC         string
	ExpiredAt    time.Time
	IsRevoked    bool
	LokasiParkir *LokasiParkir
	Petugas      *User
}

type (
	QRSessionRepository interface {
		CreateQRSession(ctx context.Context, session *QRSession) error
		GetActiveQRByPetugasID(ctx context.Context, petugasID string) (*QRSession, error)
		GetLokasiByPetugasID(ctx context.Context, petugasID string) (*QRSession, error)
		RevokeQRSession(ctx context.Context, id string) error
		GetQRSessionByPayload(ctx context.Context, payload string) (*QRSession, error)
	}

	QRSessionUsecase interface {
		CreateQRSession(ctx context.Context, petugasID string) (res *dto.QRSessionResponseDTO, err error)
		GetActiveQRByPetugasID(ctx context.Context, petugasID string) (res *dto.QRSessionResponseDTO, err error)
		GetLokasiByPetugasID(ctx context.Context, petugasID string) (res *dto.QRSessionResponseDTO, err error)
		RevokeQRSession(ctx context.Context, id string) error
		GetQRSessionByPayload(ctx context.Context, payload string) (*dto.QRSessionResponseDTO, error)
	}

	QRSessionHandler interface {
		CreateQRSession(c *gin.Context)
	}
)
