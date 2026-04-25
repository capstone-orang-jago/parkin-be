package repository

import (
	"context"
	"database/sql"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/sirupsen/logrus"
)

type transaksiParkirRepository struct {
	db *sql.DB
}

func NewTransaksiParkirRepository(db *sql.DB) *transaksiParkirRepository {
	return &transaksiParkirRepository{db: db}
}

func (repo *transaksiParkirRepository) CreateTransaksiParkir(ctx context.Context, transaksi *domain.TransaksiParkir) error {
	query := `
		INSERT INTO transaksi_parkir (kode_transaksi, pengguna_id, petugas_id, lokasi_id, jenis_kendaraan, tarif, saldo_sebelum, saldo_sesudah, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := repo.db.ExecContext(ctx, query, transaksi.KodeTransaksi, transaksi.PenggunaID, transaksi.PetugasID, transaksi.LokasiID, transaksi.JenisKendaraan, transaksi.Tarif, transaksi.SaldoSebelum, transaksi.SaldoSesudah, transaksi.Status)
	if err != nil {
		logrus.WithError(err).Error("repo.CreateTransaksiParkir: failed to create transaksi parkir")
		return err
	}
	return nil
}

func (repo *transaksiParkirRepository) GetTransaksiParkirByID(ctx context.Context, id string) (*domain.TransaksiParkir, error) {
	query := `
		SELECT id, kode_transaksi, pengguna_id, petugas_id, lokasi_id, jenis_kendaraan, tarif, saldo_sebelum, saldo_sesudah, status
		FROM transaksi_parkir
		WHERE id = $1
	`
	transaksi := &domain.TransaksiParkir{}
	err := repo.db.QueryRowContext(ctx, query, id).Scan(&transaksi.ID, &transaksi.KodeTransaksi, &transaksi.PenggunaID, &transaksi.PetugasID, &transaksi.LokasiID, &transaksi.JenisKendaraan, &transaksi.Tarif, &transaksi.SaldoSebelum, &transaksi.SaldoSesudah, &transaksi.Status)
	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetTransaksiParkirByID: no transaksi parkir found for ID")
		return nil, nil
	} else if err != nil {
		logrus.WithError(err).Error("repo.GetTransaksiParkirByID: failed to get transaksi parkir by ID")
		return nil, err
	}
	return transaksi, nil
}

func (repo *transaksiParkirRepository) GetAllTransaksiParkirByPenggunaID(ctx context.Context, penggunaID string) (transaksis []*domain.TransaksiParkir, err error) {
	query := `
		SELECT id, kode_transaksi, pengguna_id, petugas_id, lokasi_id, jenis_kendaraan, tarif, saldo_sebelum, saldo_sesudah, status
		FROM transaksi_parkir
		WHERE pengguna_id = $1
		ORDER BY created_at DESC
	`
	rows, err := repo.db.QueryContext(ctx, query, penggunaID)
	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetAllTransaksiParkirByPenggunaID: no transaksi parkir found for pengguna ID")
		return nil, nil
	} else if err != nil {
		logrus.WithError(err).Error("repo.GetAllTransaksiParkirByPenggunaID: failed to get transaksi parkir by pengguna ID")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		transaksi := &domain.TransaksiParkir{}
		err = rows.Scan(&transaksi.ID, &transaksi.KodeTransaksi, &transaksi.PenggunaID, &transaksi.PetugasID, &transaksi.LokasiID, &transaksi.JenisKendaraan, &transaksi.Tarif, &transaksi.SaldoSebelum, &transaksi.SaldoSesudah, &transaksi.Status)
		if err != nil {
			logrus.WithError(err).Error("repo.GetAllTransaksiParkirByPenggunaID: failed to scan transaksi parkir")
			return nil, err
		}
		transaksis = append(transaksis, transaksi)
	}
	return transaksis, nil
}
