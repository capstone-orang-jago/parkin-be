package domain

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type QRSession struct {
	ID         string
	LocationID string
	PetugasID  string
	Payload    string
	HMAC       string
	ExpiresAt  time.Time
	IsRevoked  bool
}

type (
	QRSessionRepository interface {
		CreateQRSession(ctx context.Context, session *QRSession) error
		GetActiveQRByPetugasID(ctx context.Context, petugasID string) (*QRSession, error)
		GetLokasiByPetugasID(ctx context.Context, petugasID string) (string, error)
		RevokeQRSession(ctx context.Context, id string) error
	}

	QRSessionUsecase interface {
		CreateQRSession(ctx context.Context, petugasID string) (res *QRSession, err error)
		GetActiveQRByPetugasID(ctx context.Context, petugasID string) (res *QRSession, err error)
		GetLokasiByPetugasID(ctx context.Context, petugasID string) (string, error)
		RevokeQRSession(ctx context.Context, id string) error
	}

	QRSessionHandler interface {
		CreateQRSession(c *gin.Context)
	}
)
