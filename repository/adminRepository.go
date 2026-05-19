package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/sirupsen/logrus"
)

type adminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) domain.AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) GetAllTransactions(ctx context.Context, q domain.AdminTransaksiQuery) ([]*domain.TransaksiParkir, int64, error) {
	baseQuery := `
		FROM transaksi_parkir t
		LEFT JOIN lokasi_parkir l ON t.lokasi_id = l.id
		LEFT JOIN users u ON t.pengguna_id = u.id
		WHERE 1=1`

	var args []interface{}
	counter := 1

	if q.DateFrom != "" {
		baseQuery += fmt.Sprintf(" AND t.created_at >= $%d", counter)
		args = append(args, q.DateFrom)
		counter++
	}
	if q.DateTo != "" {
		baseQuery += fmt.Sprintf(" AND t.created_at <= $%d", counter)
		args = append(args, q.DateTo)
		counter++
	}
	if q.LocationID != "" {
		baseQuery += fmt.Sprintf(" AND t.lokasi_id = $%d", counter)
		args = append(args, q.LocationID)
		counter++
	}
	if q.OfficerID != "" {
		baseQuery += fmt.Sprintf(" AND t.petugas_id = $%d", counter)
		args = append(args, q.OfficerID)
		counter++
	}
	if q.Vehicle != "" {
		baseQuery += fmt.Sprintf(" AND t.jenis_kendaraan = $%d", counter)
		args = append(args, q.Vehicle)
		counter++
	}
	if q.Status != "" {
		baseQuery += fmt.Sprintf(" AND t.status = $%d", counter)
		args = append(args, q.Status)
		counter++
	}

	var totalRows int64
	exec := helpers.GetExecutor(ctx, r.db)
	err := exec.QueryRowContext(ctx, "SELECT COUNT(t.id) "+baseQuery, args...).Scan(&totalRows)
	if err != nil {
		logrus.WithError(err).Error("adminRepository.GetAllTransactions: failed to count total rows")
		return nil, 0, err
	}

	offset := (q.Page - 1) * q.Limit
	dataQuery := fmt.Sprintf("SELECT t.id, t.kode_transaksi, t.jenis_kendaraan, t.tarif, t.saldo_sebelum, t.saldo_sesudah, t.status, t.created_at, l.id, l.nama_lokasi, u.id, u.name, u.email "+baseQuery+" ORDER BY t.created_at DESC LIMIT $%d OFFSET $%d", counter, counter+1)
	args = append(args, q.Limit, offset)

	rows, err := exec.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		logrus.WithError(err).Error("adminRepository.GetAllTransactions: failed to query transaction data")
		return nil, 0, err
	}
	defer rows.Close()

	var results []*domain.TransaksiParkir
	for rows.Next() {
		t := &domain.TransaksiParkir{LokasiParkir: &domain.LokasiParkir{}, Pengguna: &domain.User{}}
		err = rows.Scan(
			&t.ID, &t.KodeTransaksi, &t.JenisKendaraan, &t.Tarif, &t.SaldoSebelum, &t.SaldoSesudah, &t.Status, &t.CreatedAt,
			&t.LokasiParkir.ID, &t.LokasiParkir.NamaLokasi,
			&t.Pengguna.ID, &t.Pengguna.Name, &t.Pengguna.Email,
		)
		if err != nil {
			logrus.WithError(err).Error("adminRepository.GetAllTransactions: failed to scan transaction data")
			return nil, 0, err
		}
		results = append(results, t)
	}

	return results, totalRows, nil
}

func (r *adminRepository) GetTransactionByID(ctx context.Context, id string) (*domain.TransaksiParkir, error) {
	query := `
		SELECT t.id, t.kode_transaksi, t.jenis_kendaraan, t.tarif, t.saldo_sebelum, t.saldo_sesudah, t.status, t.created_at,
		l.id, l.nama_lokasi,
		u.id, u.name, u.email,
		p.id, p.name, p.email
		FROM transaksi_parkir t
		LEFT JOIN lokasi_parkir l ON t.lokasi_id = l.id
		LEFT JOIN users u ON t.pengguna_id = u.id
		LEFT JOIN users p ON t.petugas_id = p.id
		WHERE t.id = $1`

	t := &domain.TransaksiParkir{
		LokasiParkir: &domain.LokasiParkir{},
		Pengguna:     &domain.User{},
	}
	petugas := &domain.User{}

	var locID, locName, userID, userName, userEmail, petID, petName, petEmail sql.NullString

	exec := helpers.GetExecutor(ctx, r.db)
	err := exec.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.KodeTransaksi, &t.JenisKendaraan, &t.Tarif, &t.SaldoSebelum, &t.SaldoSesudah, &t.Status, &t.CreatedAt,
		&locID, &locName,
		&userID, &userName, &userEmail,
		&petID, &petName, &petEmail,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.WithField("id", id).Info("adminRepository.GetTransactionByID: transaction not found")
			return nil, nil
		}
		logrus.WithError(err).Error("adminRepository.GetTransactionByID: failed to query transaction data")
		return nil, err
	}

	if locID.Valid {
		t.LokasiParkir.ID = locID.String
		t.LokasiParkir.NamaLokasi = locName.String
	}
	if userID.Valid {
		t.Pengguna.ID = userID.String
		t.Pengguna.Name = userName.String
		t.Pengguna.Email = userEmail.String
	}
	if petID.Valid {
		petugas.ID = petID.String
		petugas.Name = petName.String
		petugas.Email = petEmail.String
		t.LokasiParkir.Petugas = petugas
	}

	return t, nil
}

func (r *adminRepository) GetFilterOptions(ctx context.Context) (locations []*dto.FilterOptionItem, officers []*dto.FilterOptionItem, err error) {
	exec := helpers.GetExecutor(ctx, r.db)

	rowsLoc, err := exec.QueryContext(ctx, "SELECT id, nama_lokasi FROM lokasi_parkir WHERE is_active = true")
	if err == nil {
		defer rowsLoc.Close()
		for rowsLoc.Next() {
			var loc dto.FilterOptionItem
			if err := rowsLoc.Scan(&loc.ID, &loc.Name); err == nil {
				locations = append(locations, &loc)
			}
		}
	}

	rowsOff, err := exec.QueryContext(ctx, "SELECT id, name FROM users WHERE role = 'petugas'")
	if err == nil {
		defer rowsOff.Close()
		for rowsOff.Next() {
			var off dto.FilterOptionItem
			if err := rowsOff.Scan(&off.ID, &off.Name); err == nil {
				officers = append(officers, &off)
			}
		}
	}

	return locations, officers, nil
}

func (r *adminRepository) GetAllLocations(ctx context.Context) ([]*domain.LokasiParkir, error) {
	query := `
		SELECT l.id, l.nama_lokasi, l.alamat, l.is_active, l.created_at, 
		u.id, u.name 
		FROM lokasi_parkir l
		LEFT JOIN users u ON l.petugas_id = u.id
		ORDER BY l.created_at DESC`

	exec := helpers.GetExecutor(ctx, r.db)
	rows, err := exec.QueryContext(ctx, query)
	if err != nil {
		logrus.WithError(err).Error("adminRepository.GetAllLocations: failed to query locations")
		return nil, err
	}
	defer rows.Close()

	var list []*domain.LokasiParkir
	for rows.Next() {
		loc := &domain.LokasiParkir{}
		var petID, petName sql.NullString

		err := rows.Scan(&loc.ID, &loc.NamaLokasi, &loc.Alamat, &loc.IsActive, &loc.CreatedAt, &petID, &petName)
		if err != nil {
			logrus.WithError(err).Error("adminRepository.GetAllLocations: failed to scan location data")
			return nil, err
		}

		if petID.Valid {
			loc.Petugas = &domain.User{ID: petID.String, Name: petName.String}
		}
		list = append(list, loc)
	}
	return list, nil
}

func (r *adminRepository) GetLocationStatsByID(ctx context.Context, id string) (*dto.AdminLocationDetailDTO, error) {
	query := `
		SELECT l.id, l.nama_lokasi, l.alamat, l.is_active,
		COUNT(t.id) as total_transaksi,
		COALESCE(SUM(t.tarif), 0) as total_pendapatan
		FROM lokasi_parkir l
		LEFT JOIN transaksi_parkir t ON l.id = t.lokasi_id AND t.status = 'sukses'
		WHERE l.id = $1
		GROUP BY l.id`

	exec := helpers.GetExecutor(ctx, r.db)
	stat := &dto.AdminLocationDetailDTO{}

	err := exec.QueryRowContext(ctx, query, id).Scan(
		&stat.ID, &stat.NamaLokasi, &stat.Alamat, &stat.IsActive,
		&stat.TotalTransaksi, &stat.TotalPendapatan,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.WithField("id", id).Info("adminRepository.GetLocationStatsByID: location not found")
			return nil, nil
		}
		logrus.WithError(err).Error("adminRepository.GetLocationStatsByID: failed to query location stats")
		return nil, err
	}

	return stat, nil
}
func (r *adminRepository) GetAllOfficers(ctx context.Context) ([]*domain.User, error) {
	query := `SELECT id, name, email, telephone, is_verified, created_at FROM users WHERE role = 'petugas' ORDER BY name ASC`

	exec := helpers.GetExecutor(ctx, r.db)
	rows, err := exec.QueryContext(ctx, query)
	if err != nil {
		logrus.WithError(err).Error("adminRepository.GetAllOfficers: failed to query officers")
		return nil, err
	}
	defer rows.Close()

	var list []*domain.User
	for rows.Next() {
		u := &domain.User{Role: "petugas"}
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Telephone, &u.IsVerified, &u.CreatedAt); err != nil {
			logrus.WithError(err).Error("adminRepository.GetAllOfficers: failed to scan officer data")
			return nil, err
		}
		list = append(list, u)
	}
	return list, nil
}

func (r *adminRepository) GetOfficerPerformanceByID(ctx context.Context, id string) (*dto.AdminOfficerDetailDTO, error) {
	queryUser := `SELECT id, name, email, telephone FROM users WHERE id = $1 AND role = 'petugas'`
	exec := helpers.GetExecutor(ctx, r.db)

	perf := &dto.AdminOfficerDetailDTO{}
	err := exec.QueryRowContext(ctx, queryUser, id).Scan(&perf.ID, &perf.Name, &perf.Email, &perf.Telephone)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.WithField("id", id).Info("adminRepository.GetOfficerPerformanceByID: officer not found")
			return nil, nil
		}
		logrus.WithError(err).Error("adminRepository.GetOfficerPerformanceByID: failed to query officer data")
		return nil, err
	}

	queryStats := `
		SELECT 
		COUNT(id) as total_tx,
		COUNT(DISTINCT DATE(created_at)) as total_days
		FROM transaksi_parkir 
		WHERE petugas_id = $1 AND status = 'sukses'`

	err = exec.QueryRowContext(ctx, queryStats, id).Scan(&perf.TotalTxHandled, &perf.TotalDutyDays)
	if err != nil {
		logrus.WithError(err).Error("adminRepository.GetOfficerPerformanceByID: failed to query officer stats")
		return nil, err
	}

	return perf, nil
}

func (r *adminRepository) GetRevenueReport(ctx context.Context, q domain.AdminReportQuery) (*dto.AdminReportResponseDTO, error) {
	whereClause := "WHERE t.status = 'sukses'"
	var args []interface{}
	argCount := 1

	if q.DateFrom != "" {
		whereClause += fmt.Sprintf(" AND t.created_at >= $%d", argCount)
		args = append(args, q.DateFrom)
		argCount++
	}
	if q.DateTo != "" {
		whereClause += fmt.Sprintf(" AND t.created_at <= $%d", argCount)
		args = append(args, q.DateTo)
		argCount++
	}
	if q.LocationID != "" {
		whereClause += fmt.Sprintf(" AND t.lokasi_id = $%d", argCount)
		args = append(args, q.LocationID)
		argCount++
	}
	if q.Vehicle != "" {
		whereClause += fmt.Sprintf(" AND t.jenis_kendaraan = $%d", argCount)
		args = append(args, q.Vehicle)
		argCount++
	}

	exec := helpers.GetExecutor(ctx, r.db)
	report := &dto.AdminReportResponseDTO{}

	queryGlobal := fmt.Sprintf(`
		SELECT COALESCE(SUM(t.tarif), 0), COUNT(t.id) 
		FROM transaksi_parkir t %s`, whereClause)

	err := exec.QueryRowContext(ctx, queryGlobal, args...).Scan(&report.TotalRevenue, &report.TotalTxGlobal)
	if err != nil {
		logrus.WithError(err).Error("adminRepository.GetRevenueReport: failed to query global report data")
		return nil, err
	}

	if report.TotalTxGlobal > 0 {
		report.AveragePerTx = float64(report.TotalRevenue) / float64(report.TotalTxGlobal)
	}

	queryBreakdown := fmt.Sprintf(`
		SELECT l.nama_lokasi, COUNT(t.id), COALESCE(SUM(t.tarif), 0)
		FROM transaksi_parkir t
		LEFT JOIN lokasi_parkir l ON t.lokasi_id = l.id
		%s
		GROUP BY l.nama_lokasi
		ORDER BY SUM(t.tarif) DESC`, whereClause)

	rows, err := exec.QueryContext(ctx, queryBreakdown, args...)
	if err != nil {
		logrus.WithError(err).Error("adminRepository.GetRevenueReport: failed to query breakdown data")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		br := &dto.ReportLocationBreakdown{}
		var locName sql.NullString

		if err := rows.Scan(&locName, &br.TotalTx, &br.Revenue); err == nil {
			br.LocationName = "Lokasi Tidak Diketahui"
			if locName.Valid {
				br.LocationName = locName.String
			}

			report.BreakdownBySpot = append(report.BreakdownBySpot, br)
		}
	}

	return report, nil
}
