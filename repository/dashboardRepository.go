package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/helpers"
)

type dashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) domain.DashboardRepository {
	return &dashboardRepository{db: db}
}

func (repo *dashboardRepository) GetWeeklyDashboardData(ctx context.Context, lokasiID string) (motor, mobil int, pendapatan int64, err error) {
	query := `
		SELECT 
			COUNT(CASE WHEN jenis_kendaraan = 'motor' THEN 1 END) as total_motor,
			COUNT(CASE WHEN jenis_kendaraan = 'mobil' THEN 1 END) as total_mobil,
			COALESCE(SUM(tarif), 0) as total_pendapatan
		FROM transaksi_parkir
		WHERE lokasi_id = $1 
		AND status = 'sukses' 
		AND created_at >= NOW() - INTERVAL '7 days'`

	exec := helpers.GetExecutor(ctx, repo.db)
	err = exec.QueryRowContext(ctx, query, lokasiID).Scan(&motor, &mobil, &pendapatan)
	return
}

func (repo *dashboardRepository) GetRawTransactionsForAdmin(ctx context.Context, startTime time.Time) (transaksis []*domain.TransaksiParkir, err error) {
	query := `
		SELECT t.id, t.jenis_kendaraan, t.tarif, t.status, t.created_at,
		l.id, l.nama_lokasi,
		p.id, p.name
		FROM transaksi_parkir t
		LEFT JOIN lokasi_parkir l ON t.lokasi_id = l.id
		LEFT JOIN users p ON t.petugas_id = p.id
		WHERE t.status = 'sukses' AND t.created_at >= $1
		ORDER BY t.created_at DESC`

	exec := helpers.GetExecutor(ctx, repo.db)
	rows, err := exec.QueryContext(ctx, query, startTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := &domain.TransaksiParkir{
			LokasiParkir: &domain.LokasiParkir{},
			Pengguna:     &domain.User{},
		}

		var petugasID, petugasName sql.NullString

		err = rows.Scan(
			&t.ID, &t.JenisKendaraan, &t.Tarif, &t.Status, &t.CreatedAt,
			&t.LokasiParkir.ID, &t.LokasiParkir.NamaLokasi,
			&petugasID, &petugasName,
		)
		if err != nil {
			return nil, err
		}

		if petugasID.Valid {
			t.LokasiParkir.Petugas = &domain.User{
				ID:   petugasID.String,
				Name: petugasName.String,
			}
		}

		transaksis = append(transaksis, t)
	}
	return transaksis, nil
}
