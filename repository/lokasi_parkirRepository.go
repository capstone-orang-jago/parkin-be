package repository

import (
	"context"
	"database/sql"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/sirupsen/logrus"
)

type lokasiParkirRepository struct {
	db *sql.DB
}

func NewLokasiParkirRepository(db *sql.DB) domain.LokasiParkirRepository {
	return &lokasiParkirRepository{db: db}
}

func (repo *lokasiParkirRepository) CreateLokasiParkir(ctx context.Context, lokasiParkir *domain.LokasiParkir) error {
	query := `INSERT INTO lokasi_parkir (nama_lokasi, alamat, petugas_id, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := repo.db.ExecContext(ctx, query, lokasiParkir.NamaLokasi, lokasiParkir.Alamat, lokasiParkir.PetugasID, lokasiParkir.IsActive, lokasiParkir.CreatedAt, lokasiParkir.UpdatedAt)
	if err != nil {
		logrus.WithError(err).Error("repo.CreateLokasiParkir: failed to create lokasi parkir")
		return err
	}

	return nil
}

func (repo *lokasiParkirRepository) GetAllLokasiParkir(ctx context.Context) (lokasiParkir []*domain.LokasiParkir, err error) {
	query := `SELECT l.id, l.nama_lokasi, l.alamat, l.is_active,
	p.id, p.name, p.email
	FROM lokasi_parkir l
	LEFT JOIN users p ON l.petugas_id = p.id`

	rows, err := repo.db.QueryContext(ctx, query)
	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetAllLokasiParkir: no lokasi parkir found")
		return nil, nil
	} else if err != nil {
		logrus.WithError(err).Error("repo.GetAllLokasiParkir: failed to get all lokasi parkir")
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		lokasi := &domain.LokasiParkir{
			Petugas: &domain.User{},
		}
		err = rows.Scan(&lokasi.ID, &lokasi.NamaLokasi, &lokasi.Alamat, &lokasi.IsActive, &lokasi.Petugas.ID, &lokasi.Petugas.Name, &lokasi.Petugas.Email)
		if err != nil {
			logrus.WithError(err).Error("repo.GetAllLokasiParkir: failed to scan lokasi parkir")
			return nil, err
		}
		lokasiParkir = append(lokasiParkir, lokasi)
	}

	return lokasiParkir, err
}

func (repo *lokasiParkirRepository) GetLokasiParkirByID(ctx context.Context, id string) (lokasiParkir *domain.LokasiParkir, err error) {
	query := `SELECT l.id, l.nama_lokasi, l.alamat, l.is_active,
	p.id, p.name, p.email, p.role
	FROM lokasi_parkir l
	LEFT JOIN users p ON l.petugas_id = p.id
	WHERE l.id = $1`

	lokasiParkir = &domain.LokasiParkir{
		Petugas: &domain.User{},
	}
	err = repo.db.QueryRowContext(ctx, query, id).Scan(&lokasiParkir.ID, &lokasiParkir.NamaLokasi, &lokasiParkir.Alamat, &lokasiParkir.IsActive, &lokasiParkir.Petugas.ID, &lokasiParkir.Petugas.Name, &lokasiParkir.Petugas.Email, &lokasiParkir.Petugas.Role)
	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetLokasiParkirByID: no lokasi parkir found for ID")
		return nil, nil
	} else if err != nil {
		logrus.WithError(err).Error("repo.GetLokasiParkirByID: failed to get lokasi parkir by ID")
		return nil, err
	}

	return lokasiParkir, err
}

func (repo *lokasiParkirRepository) GetLokasiParkirByPetugasID(ctx context.Context, petugasID string) (lokasiParkir *domain.LokasiParkir, err error) {
	query := `SELECT l.id, l.nama_lokasi, l.alamat, l.is_active, 
	p.id, p.name, p.email
	FROM lokasi_parkir l
	LEFT JOIN users p ON l.petugas_id = p.id
	WHERE l.petugas_id = $1`

	lokasiParkir = &domain.LokasiParkir{
		Petugas: &domain.User{},
	}

	err = repo.db.QueryRowContext(ctx, query, petugasID).Scan(&lokasiParkir.ID, &lokasiParkir.NamaLokasi, &lokasiParkir.Alamat, &lokasiParkir.IsActive, &lokasiParkir.Petugas.ID, &lokasiParkir.Petugas.Name, &lokasiParkir.Petugas.Email)
	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetLokasiParkirByPetugasID: no lokasi parkir found for PetugasID")
		return nil, nil
	} else if err != nil {
		logrus.WithError(err).Error("repo.GetLokasiParkirByPetugasID: failed to get lokasi parkir by PetugasID")
		return nil, err
	}

	return lokasiParkir, err
}

func (repo *lokasiParkirRepository) UpdateLokasiParkir(ctx context.Context, lokasiParkir *domain.LokasiParkir, id string) error {
	query := `
			UPDATE lokasi_parkir 
			SET nama_lokasi = $1, alamat = $2, petugas_id = $3, is_active = $4, updated_at = $5 
			WHERE id = $6
			`
	_, err := repo.db.ExecContext(ctx, query, lokasiParkir.NamaLokasi, lokasiParkir.Alamat, lokasiParkir.PetugasID, lokasiParkir.IsActive, lokasiParkir.UpdatedAt, id)
	if err != nil {
		logrus.WithError(err).Error("repo.UpdateLokasiParkir: failed to update lokasi parkir")
		return err
	}
	return nil
}

func (repo *lokasiParkirRepository) DeleteLokasiParkir(ctx context.Context, id string) error {
	query := `
			DELETE FROM lokasi_parkir 
			WHERE id = $1
			`
	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		logrus.WithError(err).Error("repo.DeleteLokasiParkir: failed to delete lokasi parkir")
		return err
	}
	return nil
}
