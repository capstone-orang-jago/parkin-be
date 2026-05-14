package repository

import (
	"context"
	"database/sql"

	"github.com/capstone-orang-jago/parkin-be/domain"
)

type otpCodeRepository struct {
	db *sql.DB
}

func NewOtpCodeRepository(db *sql.DB) *otpCodeRepository {
	return &otpCodeRepository{db: db}
}

func (repo *otpCodeRepository) CreateOtpCode(ctx context.Context, otpCode *domain.OtpCode) error {
	query := `
		INSERT INTO otp_codes (user_id, code, expired_at, attempts, is_used, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := repo.db.ExecContext(ctx, query, otpCode.UserId, otpCode.Code, otpCode.ExpiredAt, otpCode.Attempts, otpCode.IsUsed, otpCode.CreatedAt)
	return err
}

func (repo *otpCodeRepository) GetOtpCodeByLastActive(ctx context.Context, userId string) (*domain.OtpCode, error) {
	query := `
		SELECT id, user_id, code, expired_at, attempts, is_used, created_at
		FROM otp_codes
		WHERE user_id = $1 AND is_used = false
		ORDER BY created_at DESC
		LIMIT 1
	`
	otpCode := &domain.OtpCode{}
	err := repo.db.QueryRowContext(ctx, query, userId).Scan(&otpCode.ID, &otpCode.UserId, &otpCode.Code, &otpCode.ExpiredAt, &otpCode.Attempts, &otpCode.IsUsed, &otpCode.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return otpCode, err
}

func (repo *otpCodeRepository) UpdateOtpCode(ctx context.Context, otp *domain.OtpCode) error {
	query := `UPDATE otp_codes SET attempts = $1, is_used = $2 WHERE id = $3`
	_, err := repo.db.ExecContext(ctx, query, otp.Attempts, otp.IsUsed, otp.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *otpCodeRepository) DeleteOtpCode(ctx context.Context, id string) error {
	query := `DELETE FROM otp_codes WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
