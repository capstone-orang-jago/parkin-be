package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/capstone-orang-jago/parkin-be/config"
	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type authUsecase struct {
	authRepo domain.AuthRepository
	userRepo domain.UserRepository
}

func NewAuthUsecase(authRepo domain.AuthRepository, userRepo domain.UserRepository) domain.AuthUsecase {
	return &authUsecase{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

func (usecase *authUsecase) Login(ctx context.Context, req dto.LoginRequestDTO) (*dto.LoginResponseDTO, error) {
	if req.Email == "" {
		return nil, errors.New("email are required")
	}
	if req.Password == "" {
		return nil, errors.New("password are required")
	}

	user, err := usecase.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	_, err = helpers.ComparePass(user.Password, req.Password)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	refreshToken, err := helpers.GenerateRefreshToken()
	if err != nil {
		logrus.WithError(err).Error("usecase.Login: failed to generate refresh token")
		return nil, errors.New("failed to generate refresh token")
	}

	accessToken, err := helpers.GenerateAccessToken(user.ID, config.AppConfig.JWTSecret, time.Hour*2)
	if err != nil {
		logrus.WithError(err).Error("usecase.Login: failed to generate access token")
		return nil, errors.New("failed to generate access token")
	}

	tokenId := uuid.NewString()

	auth := &domain.Auth{
		ID:        tokenId,
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}

	if err = usecase.authRepo.CreateRefreshToken(ctx, auth); err != nil {
		logrus.WithError(err).Error("usecase.Login: failed to save refresh token to database")
		return nil, errors.New("failed to save refresh token to database")
	}

	res := &dto.LoginResponseDTO{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (usecase *authUsecase) RefreshToken(ctx context.Context, req dto.RefreshTokenRequestDTO) (*domain.Auth, error) {
	if req.RefreshToken == "" {
		return nil, errors.New("refresh token are required")
	}

	auth, err := usecase.authRepo.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		logrus.WithError(err).Error("usecase.RefreshToken: invalid refresh token")
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().After(auth.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	accessToken, err := helpers.GenerateAccessToken(
		auth.UserID,
		config.AppConfig.JWTSecret,
		time.Hour*2,
	)
	if err != nil {
		return nil, err
	}

	data := &domain.Auth{
		ID:        auth.ID,
		UserID:    auth.UserID,
		Token:     accessToken,
		ExpiresAt: time.Now().Add(time.Hour * 2),
	}

	return data, nil
}

func (usecase *authUsecase) Logout(ctx context.Context, token string) error {
	if token == "" {
		return errors.New("refresh token are required")
	}

	err := usecase.authRepo.DeleteRefreshToken(ctx, token)
	if err != nil {
		logrus.WithError(err).Error("usecase.Logout: failed to delete refresh token from database")
		return errors.New("failed to delete refresh token from database")
	}

	return nil
}
