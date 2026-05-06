package domain

import (
	"context"
	"time"

	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	ID               string
	UserID           string
	RefreshTokenHash string
	ExpiredAt        time.Time
}

type (
	AuthRepository interface {
		CreateRefreshToken(ctx context.Context, auth *Auth) error
		GetRefreshToken(ctx context.Context, token string) (*Auth, error)
		DeleteRefreshToken(ctx context.Context, token string) error
	}

	AuthUsecase interface {
		Login(ctx context.Context, req dto.LoginRequestDTO) (res *dto.LoginResponseDTO, err error)
		RefreshToken(ctx context.Context, req dto.RefreshTokenRequestDTO) (res *dto.RefreshTokenResponseDTO, err error)
		Logout(ctx context.Context, userId, token string) error
	}

	AuthHandler interface {
		Register(c *gin.Context)
		Login(c *gin.Context)
		RefreshToken(c *gin.Context)
		Logout(c *gin.Context)
	}
)
