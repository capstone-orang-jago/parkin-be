package repository

import (
	"context"
	"database/sql"

	"github.com/capstone-orang-jago/parkin-be/domain"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) domain.AuthRepository {
	return &authRepository{db: db}
}

func (repo *authRepository) CreateRefreshToken(ctx context.Context, auth *domain.Auth) error {
	query := "INSERT INTO refresh_tokens (id, user_id, token, expires_at) VALUES ($1, $2, $3, $4)"
	_, err := repo.db.ExecContext(ctx, query, auth.ID, auth.UserID, auth.Token, auth.ExpiresAt)
	return err
}

func (repo *authRepository) GetRefreshToken(ctx context.Context, token string) (*domain.Auth, error) {
	query := "SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token = $1"
	auth := &domain.Auth{}
	err := repo.db.QueryRowContext(ctx, query, token).Scan(&auth.ID, &auth.UserID, &auth.Token, &auth.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (repo *authRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	query := "DELETE FROM refresh_tokens WHERE token = $1"
	_, err := repo.db.ExecContext(ctx, query, token)
	return err
}
