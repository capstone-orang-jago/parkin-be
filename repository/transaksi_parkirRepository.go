package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/sirupsen/logrus"
)

type transaksiParkirRepository struct {
	db *sql.DB
}

func NewTransaksiParkirRepository(db *sql.DB) *transaksiParkirRepository {
	return &transaksiParkirRepository{db: db}
}

func (repo *transaksiParkirRepository) CreateTransaksiParkir(ctx context.Context, transaksi *domain.TransaksiParkir) error {
	queryInsert := `
        INSERT INTO transaksi_parkir (kode_transaksi, pengguna_id, petugas_id, lokasi_id, jenis_kendaraan, tarif, saldo_sebelum, saldo_sesudah, status)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id
    `

	exec := helpers.GetExecutor(ctx, repo.db)
	err := exec.QueryRowContext(ctx, queryInsert, transaksi.KodeTransaksi, transaksi.PenggunaID, transaksi.PetugasID, transaksi.LokasiID, transaksi.JenisKendaraan, transaksi.Tarif, transaksi.SaldoSebelum, transaksi.SaldoSesudah, transaksi.Status).Scan(&transaksi.ID)
	if err != nil {
		logrus.WithError(err).Error("repo.CreateTransaksiParkir: failed to create transaksi parkir")
		return err
	}

	querySelect := `
        SELECT 
        t.id, t.kode_transaksi, t.jenis_kendaraan, t.tarif, t.saldo_sebelum, t.saldo_sesudah, t.status, t.created_at,
        l.id, l.nama_lokasi, l.alamat, l.is_active,
        u.id, u.name, u.email,
        p.id, p.name, p.email
        FROM transaksi_parkir t
        LEFT JOIN lokasi_parkir l ON t.lokasi_id = l.id
        LEFT JOIN users u ON t.pengguna_id = u.id
        LEFT JOIN users p ON t.petugas_id = p.id
        WHERE t.id = $1
    `

	transaksi.LokasiParkir = &domain.LokasiParkir{
		Petugas: &domain.User{},
	}
	transaksi.Pengguna = &domain.User{}

	err = exec.QueryRowContext(ctx, querySelect, transaksi.ID).Scan(
		&transaksi.ID, &transaksi.KodeTransaksi, &transaksi.JenisKendaraan, &transaksi.Tarif, &transaksi.SaldoSebelum, &transaksi.SaldoSesudah, &transaksi.Status, &transaksi.CreatedAt,
		&transaksi.LokasiParkir.ID, &transaksi.LokasiParkir.NamaLokasi, &transaksi.LokasiParkir.Alamat, &transaksi.LokasiParkir.IsActive,
		&transaksi.Pengguna.ID, &transaksi.Pengguna.Name, &transaksi.Pengguna.Email,
		&transaksi.LokasiParkir.Petugas.ID, &transaksi.LokasiParkir.Petugas.Name, &transaksi.LokasiParkir.Petugas.Email,
	)
	if err != nil {
		logrus.WithError(err).Error("repo.CreateTransaksiParkir: failed to re-fetch transaction data")
		return err
	}

	return nil
}

func (repo *transaksiParkirRepository) GetTransaksiParkirByID(ctx context.Context, id string) (*domain.TransaksiParkir, error) {
	query := `
        SELECT t.id, t.kode_transaksi, t.jenis_kendaraan, t.tarif, t.saldo_sebelum, t.saldo_sesudah, t.status, t.created_at,
        l.id, l.nama_lokasi, l.alamat, l.is_active,
        u.id, u.name, u.email,
        p.id, p.name, p.email
        FROM transaksi_parkir t
        LEFT JOIN lokasi_parkir l ON t.lokasi_id = l.id
        LEFT JOIN users u ON t.pengguna_id = u.id
        LEFT JOIN users p ON t.petugas_id = p.id
        WHERE t.id = $1
    `
	exec := helpers.GetExecutor(ctx, repo.db)

	transaksi := &domain.TransaksiParkir{
		Pengguna: &domain.User{},
		LokasiParkir: &domain.LokasiParkir{
			Petugas: &domain.User{},
		},
	}

	err := exec.QueryRowContext(ctx, query, id).Scan(
		&transaksi.ID, &transaksi.KodeTransaksi, &transaksi.JenisKendaraan, &transaksi.Tarif, &transaksi.SaldoSebelum, &transaksi.SaldoSesudah, &transaksi.Status, &transaksi.CreatedAt,
		&transaksi.LokasiParkir.ID, &transaksi.LokasiParkir.NamaLokasi, &transaksi.LokasiParkir.Alamat, &transaksi.LokasiParkir.IsActive,
		&transaksi.Pengguna.ID, &transaksi.Pengguna.Name, &transaksi.Pengguna.Email,
		&transaksi.LokasiParkir.Petugas.ID, &transaksi.LokasiParkir.Petugas.Name, &transaksi.LokasiParkir.Petugas.Email,
	)

	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetTransaksiParkirByID: no transaksi parkir found for ID")
		return nil, nil
	} else if err != nil {
		logrus.WithError(err).Error("repo.GetTransaksiParkirByID: failed to get transaksi parkir by ID")
		return nil, err
	}
	return transaksi, nil
}

func (repo *transaksiParkirRepository) GetAllTransaksiParkirByPenggunaID(ctx context.Context, penggunaID string, page, limit int, filter string) (transaksis []*domain.TransaksiParkir, err error) {

	offset := (page - 1) * limit

	query := `
        SELECT t.id, t.kode_transaksi, t.jenis_kendaraan, t.tarif, t.saldo_sebelum, t.saldo_sesudah, t.status, t.created_at,
        l.id, l.nama_lokasi, l.alamat, l.is_active,
        u.id, u.name, u.email,
        p.id, p.name, p.email
        FROM transaksi_parkir t
        LEFT JOIN lokasi_parkir l ON t.lokasi_id = l.id
        LEFT JOIN users u ON t.pengguna_id = u.id
        LEFT JOIN users p ON t.petugas_id = p.id
        WHERE t.pengguna_id = $1
    `

	args := []interface{}{penggunaID}
	argCount := 2

	switch filter {
	case "minggu":
		query += ` AND t.created_at >= NOW() - INTERVAL '7 days'`
	case "bulan":
		query += ` AND t.created_at >= NOW() - INTERVAL '30 days'`
	}

	query += fmt.Sprintf(" ORDER BY t.created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	exec := helpers.GetExecutor(ctx, repo.db)
	rows, err := exec.QueryContext(ctx, query, args...)
	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetAllTransaksiParkirByPenggunaID: no transaksi parkir found for pengguna ID")
		return nil, nil
	} else if err != nil {
		logrus.WithError(err).Error("repo.GetAllTransaksiParkirByPenggunaID: failed to get transaksi parkir by pengguna ID")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		transaksi := &domain.TransaksiParkir{
			Pengguna: &domain.User{},
			LokasiParkir: &domain.LokasiParkir{
				Petugas: &domain.User{},
			},
		}
		err = rows.Scan(
			&transaksi.ID, &transaksi.KodeTransaksi, &transaksi.JenisKendaraan, &transaksi.Tarif, &transaksi.SaldoSebelum, &transaksi.SaldoSesudah, &transaksi.Status, &transaksi.CreatedAt,
			&transaksi.LokasiParkir.ID, &transaksi.LokasiParkir.NamaLokasi, &transaksi.LokasiParkir.Alamat, &transaksi.LokasiParkir.IsActive,
			&transaksi.Pengguna.ID, &transaksi.Pengguna.Name, &transaksi.Pengguna.Email,
			&transaksi.LokasiParkir.Petugas.ID, &transaksi.LokasiParkir.Petugas.Name, &transaksi.LokasiParkir.Petugas.Email,
		)
		if err != nil {
			logrus.WithError(err).Error("repo.GetAllTransaksiParkirByPenggunaID: failed to scan transaksi parkir")
			return nil, err
		}
		transaksis = append(transaksis, transaksi)
	}
	return transaksis, nil
}
