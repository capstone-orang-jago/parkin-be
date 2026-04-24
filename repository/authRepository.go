package repository

import (
	"context"
	"database/sql"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/sirupsen/logrus"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) domain.AuthRepository {
	return &authRepository{db: db}
}

func (repo *authRepository) CreateRefreshToken(ctx context.Context, auth *domain.Auth) error {
	query := "INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)"
	_, err := repo.db.ExecContext(ctx, query, auth.UserID, auth.RefreshTokenHash, auth.ExpiresAt)
	if err != nil {
		logrus.WithError(err).Error("repo.CreateRefreshToken: failed to create refresh token")
	}
	return err
}

func (repo *authRepository) GetRefreshToken(ctx context.Context, token string) (*domain.Auth, error) {
	query := "SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token = $1"
	auth := &domain.Auth{}
	err := repo.db.QueryRowContext(ctx, query, token).Scan(&auth.ID, &auth.UserID, &auth.RefreshTokenHash, &auth.ExpiresAt)
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
	_, err := repo.db.ExecContext(ctx, query, token)
	if err != nil {
		logrus.WithError(err).Error("repo.DeleteRefreshToken: failed to delete refresh token")
	}
	return err
}
