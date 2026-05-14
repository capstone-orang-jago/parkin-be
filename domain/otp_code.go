package domain

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type (
	OtpCode struct {
		ID        string
		UserId    string
		Code      string
		ExpiredAt time.Time
		Attempts  int
		IsUsed    bool
		CreatedAt time.Time
	}

	OtpCodeRepository interface {
		CreateOtpCode(ctx context.Context, otpCode *OtpCode) error
		GetOtpCodeByLastActive(ctx context.Context, userId string) (*OtpCode, error)
		UpdateOtpCode(ctx context.Context, otpCode *OtpCode) error
		DeleteOtpCode(ctx context.Context, id string) error
	}

	OtpCodeUsecase interface {
		SendOtpCode(ctx context.Context, telephone string) error
		VerifyOtpCode(ctx context.Context, telephone string, code string) (bool, error)
	}

	OtpCodeHandler interface {
		ResendOtpCode(c *gin.Context)
	}
)
