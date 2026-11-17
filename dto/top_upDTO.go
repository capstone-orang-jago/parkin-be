package dto

import (
	"time"
)

type TopUpRequestDTO struct {
	Nominal          int64  `json:"nominal" binding:"required"`
	MetodePembayaran string `json:"metode_pembayaran" binding:"required"`
}

type TopUpResponseDTO struct {
	ID                    string
	User                  *UserResponseDTO
	Nominal               int64
	MetodePembayaran      string
	MidtransOrderID       string
	MidtransTransactionID string
	Status                string
	SnapToken             string
	SnapRedirectURL       string
	PaidAt                time.Time
}
