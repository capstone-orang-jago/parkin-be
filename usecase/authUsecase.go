package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/capstone-orang-jago/parkin-be/config"
	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/capstone-orang-jago/parkin-be/helpers"
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

	if user == nil {
		logrus.Info("usecase.Login: user not found for email")
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

	hashedToken := helpers.HashToken(refreshToken)

	accessToken, err := helpers.GenerateAccessToken(
		user.ID,
		user.Role,
		user.Name,
		config.AppConfig.JWTSecret,
		time.Hour*2,
	)
	if err != nil {
		logrus.WithError(err).Error("usecase.Login: failed to generate access token")
		return nil, errors.New("failed to generate access token")
	}

	auth := &domain.Auth{
		UserID:           user.ID,
		RefreshTokenHash: hashedToken,
		ExpiredAt:        time.Now().Add(time.Hour * 24 * 30),
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

func (usecase *authUsecase) RefreshToken(ctx context.Context, req dto.RefreshTokenRequestDTO) (*dto.RefreshTokenResponseDTO, error) {
	if req.RefreshToken == "" {
		return nil, errors.New("refresh token are required")
	}

	oldHashedToken := helpers.HashToken(req.RefreshToken)

	auth, err := usecase.authRepo.GetRefreshToken(ctx, oldHashedToken)
	if err != nil {
		logrus.WithError(err).Error("usecase.RefreshToken: invalid refresh token")
		return nil, errors.New("invalid refresh token")
	}

	if auth == nil {
		logrus.Error("usecase.RefreshToken: refresh token not found")
		return nil, errors.New("invalid refresh token")
	}

	user, err := usecase.userRepo.GetUserByID(ctx, auth.UserID)
	if err != nil {
		logrus.WithError(err).Error("usecase.RefreshToken: failed to get user")
		return nil, errors.New("failed to get user")
	}

	if user == nil {
		logrus.Error("usecase.RefreshToken: user not found")
		return nil, errors.New("user not found")
	}

	if time.Now().After(auth.ExpiredAt) {
		logrus.Error("usecase.RefreshToken: refresh token expired")
		return nil, errors.New("refresh token expired")
	}

	newRefreshToken, err := helpers.GenerateRefreshToken()
	if err != nil {
		logrus.WithError(err).Error("usecase.RefreshToken: failed to generate new refresh token")
		return nil, errors.New("failed to generate new refresh token")
	}

	newAccessToken, err := helpers.GenerateAccessToken(
		auth.UserID,
		user.Role,
		user.Name,
		config.AppConfig.JWTSecret,
		time.Hour*2,
	)
	if err != nil {
		logrus.WithError(err).Error("usecase.RefreshToken: failed to generate new access token")
		return nil, errors.New("failed to generate new access token")
	}

	newHashedToken := helpers.HashToken(newRefreshToken)

	newAuth := &domain.Auth{
		UserID:           auth.UserID,
		RefreshTokenHash: newHashedToken,
		ExpiredAt:        time.Now().Add(time.Hour * 24 * 30),
	}

	err = usecase.authRepo.CreateRefreshToken(ctx, newAuth)
	if err != nil {
		return nil, errors.New("failed to save new refresh token to database")
	}

	err = usecase.authRepo.DeleteRefreshToken(ctx, oldHashedToken)
	if err != nil {
		return nil, errors.New("failed to delete old refresh token from database")
	}

	data := &dto.RefreshTokenResponseDTO{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	return data, nil
}

func (usecase *authUsecase) Logout(ctx context.Context, userId, token string) error {
	if token == "" {
		return errors.New("refresh token are required")
	}

	hashedToken := helpers.HashToken(token)

	auth, err := usecase.authRepo.GetRefreshToken(ctx, hashedToken)
	if err != nil {
		return err
	}

	if auth == nil {
		return errors.New("refresh token not found")
	}

	if auth.UserID != userId {
		return errors.New("unauthorized")
	}

	err = usecase.authRepo.DeleteRefreshToken(ctx, hashedToken)
	if err != nil {
		logrus.WithError(err).Error("usecase.Logout: failed to delete refresh token from database")
		return errors.New("failed to delete refresh token from database")
	}

	return nil
}
