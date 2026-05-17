package domain

import (
	"context"
	"encoding/json"
	"time"
)

type (
	TopUp struct {
		ID                    string
		UserID                string
		Nominal               int64
		MetodePembayaran      string
		MidtransOrderID       string
		MidtransTransactionID string
		Status                string
		SnapToken             string
		SnapRedirectURL       string
		RawCallBack           json.RawMessage
		PaidAt                time.Time
		CreatedAt             time.Time
		UpdatedAt             time.Time
		User                  *User
	}

	TopUpRepository interface {
		CreateTopUp(ctx context.Context, topUp *TopUp) error
		GetTopUpByMidtransOrderID(ctx context.Context, midtransOrderID string) (*TopUp, error)
		UpdateTopUp(ctx context.Context, topUp *TopUp) error
	}
)
