package repository

import (
	"context"
	"database/sql"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/helpers"
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
	queryInsert := `INSERT INTO qr_sessions (lokasi_id, petugas_id, payload, hmac, expired_at, is_revoked) VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`
	exec := helpers.GetExecutor(ctx, repo.db)
	err := exec.QueryRowContext(ctx, queryInsert, session.LokasiID, session.PetugasID, session.Payload, session.HMAC, session.ExpiredAt, session.IsRevoked).Scan(&session.ID)
	if err != nil {
		logrus.WithError(err).Error("repo.CreateQRSession: failed to create QR session")
		return err
	}

	querySelect := `
		SELECT q.id, q.payload, q.hmac, q.expired_at, q.is_revoked,
		p.id, p.name, p.email,
		l.id, l.nama_lokasi, l.alamat, l.is_active
		FROM qr_sessions q
		LEFT JOIN users p ON q.petugas_id = p.id
		LEFT JOIN lokasi_parkir l ON q.lokasi_id = l.id
		WHERE q.id = $1`

	session.Petugas = &domain.User{}
	session.LokasiParkir = &domain.LokasiParkir{}
	err = exec.QueryRowContext(ctx, querySelect, session.ID).Scan(
		&session.ID, &session.Payload, &session.HMAC, &session.ExpiredAt, &session.IsRevoked,
		&session.Petugas.ID, &session.Petugas.Name, &session.Petugas.Email,
		&session.LokasiParkir.ID, &session.LokasiParkir.NamaLokasi, &session.LokasiParkir.Alamat, &session.LokasiParkir.IsActive,
	)
	if err != nil {
		logrus.WithError(err).Error("repo.CreateQRSession: failed to re-fetch QR session with joins")
		return err
	}
	return nil
}

func (repo *qrSessionRepository) GetActiveQRByPetugasID(ctx context.Context, petugasID string) (*domain.QRSession, error) {
	query := `SELECT q.id, q.payload, q.hmac, q.expired_at, q.is_revoked,
    p.id, p.name, p.email,
    l.id, l.nama_lokasi, l.alamat, l.is_active
    FROM qr_sessions q
    LEFT JOIN users p on q.petugas_id = p.id
    LEFT JOIN lokasi_parkir l on q.lokasi_id = l.id
    WHERE q.petugas_id = $1 
    AND q.is_revoked = false 
    AND q.expired_at > NOW()  
    FOR UPDATE OF q`

	exec := helpers.GetExecutor(ctx, repo.db)
	row := exec.QueryRowContext(ctx, query, petugasID)
	qrSession := &domain.QRSession{
		Petugas:      &domain.User{},
		LokasiParkir: &domain.LokasiParkir{},
	}
	err := row.Scan(&qrSession.ID, &qrSession.Payload, &qrSession.HMAC, &qrSession.ExpiredAt, &qrSession.IsRevoked, &qrSession.Petugas.ID, &qrSession.Petugas.Name, &qrSession.Petugas.Email, &qrSession.LokasiParkir.ID, &qrSession.LokasiParkir.NamaLokasi, &qrSession.LokasiParkir.Alamat, &qrSession.LokasiParkir.IsActive)
	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetActiveQRByPetugasID: no active QR session found for petugas ID")
		return nil, nil
	} else if err != nil {
		logrus.WithError(err).Error("repo.GetActiveQRByPetugasID: failed to get active QR session by petugas ID")
		return nil, err
	}
	return qrSession, nil
}

func (repo *qrSessionRepository) RevokeQRSession(ctx context.Context, id string) error {
	query := `UPDATE qr_sessions SET is_revoked = true WHERE id = $1`
	exec := helpers.GetExecutor(ctx, repo.db)
	_, err := exec.ExecContext(ctx, query, id)
	if err != nil {
		logrus.WithError(err).Error("repo.RevokeQRSession: failed to revoke QR session")
	}
	return err
}

func (repo *qrSessionRepository) GetQRSessionByPayload(ctx context.Context, payload string) (*domain.QRSession, error) {
	query := `SELECT q.id, q.payload, q.hmac, q.expired_at, q.is_revoked,
	p.id, p.name, p.email,
	l.id, l.nama_lokasi, l.alamat, l.is_active
	FROM qr_sessions q 
	LEFT JOIN users p ON q.petugas_id = p.id
	LEFT JOIN lokasi_parkir l ON q.lokasi_id = l.id
	WHERE q.payload = $1 
	AND q.is_revoked = false 
	AND q.expired_at > NOW()`
	exec := helpers.GetExecutor(ctx, repo.db)
	row := exec.QueryRowContext(ctx, query, payload)
	qrSession := &domain.QRSession{
		Petugas:      &domain.User{},
		LokasiParkir: &domain.LokasiParkir{},
	}
	err := row.Scan(&qrSession.ID, &qrSession.Payload, &qrSession.HMAC, &qrSession.ExpiredAt, &qrSession.IsRevoked, &qrSession.Petugas.ID, &qrSession.Petugas.Name, &qrSession.Petugas.Email, &qrSession.LokasiParkir.ID, &qrSession.LokasiParkir.NamaLokasi, &qrSession.LokasiParkir.Alamat, &qrSession.LokasiParkir.IsActive)
	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetQRSessionByPayload: no QR session found for payload")
		return nil, nil
	} else if err != nil {
		logrus.WithError(err).Error("repo.GetQRSessionByPayload: failed to get QR session by payload")
		return nil, err
	}
	return qrSession, nil
}
