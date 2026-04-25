package repository

import (
	"context"
	"database/sql"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/sirupsen/logrus"
)

type qrSessionRepository struct {
	db *sql.DB
}

func NewQRSessionRepository(db *sql.DB) *qrSessionRepository {
	return &qrSessionRepository{
		db: db,
	}
}

func (repo *qrSessionRepository) CreateQRSession(ctx context.Context, session *domain.QRSession) error {
	query := `INSERT INTO qr_sessions (lokasi_id, petugas_id, payload, hmac, expires_at, is_revoked) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := repo.db.ExecContext(ctx, query, session.LocationID, session.PetugasID, session.Payload, session.HMAC, session.ExpiresAt, session.IsRevoked)
	if err != nil {
		logrus.WithError(err).Error("repo.CreateQRSession: failed to create QR session")
		return err
	}
	return nil
}

func (repo *qrSessionRepository) GetActiveQRByPetugasID(ctx context.Context, petugasID string) (*domain.QRSession, error) {
	query := `SELECT id, lokasi_id, petugas_id, payload, hmac, expires_at, is_revoked FROM qr_sessions WHERE petugas_id = $1 AND is_revoked = false AND expires_at > NOW()`
	row := repo.db.QueryRowContext(ctx, query, petugasID)
	var session domain.QRSession
	err := row.Scan(&session.ID, &session.LocationID, &session.PetugasID, &session.Payload, &session.HMAC, &session.ExpiresAt, &session.IsRevoked)
	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetActiveQRByPetugasID: no active QR session found for petugas ID")
		return nil, nil
	} else if err != nil {
		logrus.WithError(err).Error("repo.GetActiveQRByPetugasID: failed to get active QR session by petugas ID")
		return nil, err
	}
	return &session, nil
}

func (repo *qrSessionRepository) GetLokasiByPetugasID(ctx context.Context, petugasID string) (string, error) {
	query := `SELECT lokasi_id FROM qr_sessions WHERE petugas_id = $1 AND is_revoked = false AND expires_at > NOW()`
	row := repo.db.QueryRowContext(ctx, query, petugasID)

	var locationID string
	err := row.Scan(&locationID)
	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetLokasiByPetugasID: no lokasi found for petugas ID")
		return "", nil
	} else if err != nil {
		logrus.WithError(err).Error("repo.GetLokasiByPetugasID: failed to get lokasi by petugas ID")
		return "", err
	}
	return locationID, nil
}

func (repo *qrSessionRepository) RevokeQRSession(ctx context.Context, id string) error {
	query := `UPDATE qr_sessions SET is_revoked = true WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		logrus.WithError(err).Error("repo.RevokeQRSession: failed to revoke QR session")
	}
	return err
}

func (repo *qrSessionRepository) GetQRSessionByPayload(ctx context.Context, payload string) (*domain.QRSession, error) {
	query := `SELECT id, lokasi_id, petugas_id, payload, hmac, expires_at, is_revoked FROM qr_sessions WHERE payload = $1 AND is_revoked = false AND expires_at > NOW()`
	row := repo.db.QueryRowContext(ctx, query, payload)
	var session domain.QRSession
	err := row.Scan(&session.ID, &session.LocationID, &session.PetugasID, &session.Payload, &session.HMAC, &session.ExpiresAt, &session.IsRevoked)
	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetQRSessionByPayload: no QR session found for payload")
		return nil, nil
	} else if err != nil {
		logrus.WithError(err).Error("repo.GetQRSessionByPayload: failed to get QR session by payload")
		return nil, err
	}
	return &session, nil
}
