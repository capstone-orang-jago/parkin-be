package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/sirupsen/logrus"
)

type topUpRepository struct {
	db *sql.DB
}

func NewTopUpRepository(db *sql.DB) *topUpRepository {
	return &topUpRepository{
		db: db,
	}
}

func (repo *topUpRepository) CreateTopUp(ctx context.Context, topUp *domain.TopUp) error {
	query := `INSERT INTO topup_saldo (pengguna_id, nominal, metode_pembayaran, midtrans_order_id, midtrans_transaction_id, status, snap_token, snap_redirect_url, raw_callback, paid_at, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW()) RETURNING id`

	exec := helpers.GetExecutor(ctx, repo.db)
	err := exec.QueryRowContext(ctx, query,
		topUp.UserID,
		topUp.Nominal,
		topUp.MetodePembayaran,
		topUp.MidtransOrderID,
		topUp.MidtransTransactionID,
		topUp.Status,
		topUp.SnapToken,
		topUp.SnapRedirectURL,
		topUp.RawCallBack,
		topUp.PaidAt,
	).Scan(&topUp.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *topUpRepository) GetTopUpByMidtransOrderID(ctx context.Context, midtransOrderID string) (*domain.TopUp, error) {
	query := `SELECT tu.id, tu.nominal, tu.metode_pembayaran, tu.midtrans_order_id, tu.midtrans_transaction_id, tu.status, tu.snap_token, tu.snap_redirect_url, tu.raw_callback, tu.paid_at,
	u.id, u.name, u.email, u.saldo
	FROM topup_saldo tu
	LEFT JOIN users u ON tu.pengguna_id = u.id
	WHERE tu.midtrans_order_id = $1`

	exec := helpers.GetExecutor(ctx, repo.db)
	row := exec.QueryRowContext(ctx, query, midtransOrderID)

	return scanTopUp(row)
}

func (repo *topUpRepository) GetTopUpByMidtransOrderIDForUpdate(ctx context.Context, midtransOrderID string) (*domain.TopUp, error) {
	query := `SELECT tu.id, tu.nominal, tu.metode_pembayaran, tu.midtrans_order_id, tu.midtrans_transaction_id, tu.status, tu.snap_token, tu.snap_redirect_url, tu.raw_callback, tu.paid_at,
	u.id, u.name, u.email, u.saldo
	FROM topup_saldo tu
	LEFT JOIN users u ON tu.pengguna_id = u.id
	WHERE tu.midtrans_order_id = $1
	FOR UPDATE OF tu`

	exec := helpers.GetExecutor(ctx, repo.db)
	row := exec.QueryRowContext(ctx, query, midtransOrderID)

	return scanTopUp(row)
}

func scanTopUp(row *sql.Row) (*domain.TopUp, error) {

	topUp := &domain.TopUp{}
	topUp.User = &domain.User{}

	var rawCallbackTmp []byte
	var paidAtTmp sql.NullTime

	err := row.Scan(
		&topUp.ID,
		&topUp.Nominal,
		&topUp.MetodePembayaran,
		&topUp.MidtransOrderID,
		&topUp.MidtransTransactionID,
		&topUp.Status,
		&topUp.SnapToken,
		&topUp.SnapRedirectURL,
		&rawCallbackTmp,
		&paidAtTmp,
		&topUp.User.ID,
		&topUp.User.Name,
		&topUp.User.Email,
		&topUp.User.Saldo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	topUp.RawCallBack = json.RawMessage(rawCallbackTmp)
	if paidAtTmp.Valid {
		topUp.PaidAt = paidAtTmp.Time
	}

	return topUp, nil
}

func (repo *topUpRepository) UpdateTopUp(ctx context.Context, topUp *domain.TopUp) error {
	query := `UPDATE topup_saldo SET nominal = $1, metode_pembayaran = $2, midtrans_order_id = $3, midtrans_transaction_id = $4, status = $5, snap_token = $6, snap_redirect_url = $7, raw_callback = $8, paid_at = $9, updated_at = NOW() WHERE id = $10`
	exec := helpers.GetExecutor(ctx, repo.db)
	_, err := exec.ExecContext(ctx, query,
		topUp.Nominal,
		topUp.MetodePembayaran,
		topUp.MidtransOrderID,
		topUp.MidtransTransactionID,
		topUp.Status,
		topUp.SnapToken,
		topUp.SnapRedirectURL,
		topUp.RawCallBack,
		topUp.PaidAt,
		topUp.ID,
	)
	if err != nil {
		logrus.WithError(err).Error("repo.UpdateTopUp: failed to update top-up record")
		return err
	}
	return nil
}
