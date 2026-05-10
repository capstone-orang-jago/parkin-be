package repository

import (
	"context"
	"database/sql"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/sirupsen/logrus"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) domain.AuthRepository {
	return &authRepository{db: db}
}

func (repo *authRepository) CreateRefreshToken(ctx context.Context, auth *domain.Auth) error {
	query := "INSERT INTO refresh_tokens (user_id, token, expired_at) VALUES ($1, $2, $3)"

	exec := helpers.GetExecutor(ctx, repo.db)
	_, err := exec.ExecContext(ctx, query, auth.UserID, auth.RefreshTokenHash, auth.ExpiredAt)
	if err != nil {
		logrus.WithError(err).Error("repo.CreateRefreshToken: failed to create refresh token")
	}
	return err
}

func (repo *authRepository) GetRefreshToken(ctx context.Context, token string) (*domain.Auth, error) {
	query := "SELECT id, user_id, token, expired_at FROM refresh_tokens WHERE token = $1"
	auth := &domain.Auth{}
	exec := helpers.GetExecutor(ctx, repo.db)
	err := exec.QueryRowContext(ctx, query, token).Scan(&auth.ID, &auth.UserID, &auth.RefreshTokenHash, &auth.ExpiredAt)
	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetRefreshToken: no refresh token found")
		return nil, nil
	} else if err != nil {
		logrus.WithError(err).Error("repo.GetRefreshToken: failed to get refresh token")
		return nil, err
	}
	return auth, nil
}

func (repo *authRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	query := "DELETE FROM refresh_tokens WHERE token = $1"

	exec := helpers.GetExecutor(ctx, repo.db)
	_, err := exec.ExecContext(ctx, query, token)
	if err != nil {
		logrus.WithError(err).Error("repo.DeleteRefreshToken: failed to delete refresh token")
	}
	return err
}
