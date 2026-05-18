package usecase

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/sirupsen/logrus"
)

type midtransUsecase struct {
	userRepo          domain.UserRepository
	topUpRepo         domain.TopUpRepository
	transactor        helpers.Transactor
	midtransServerKey string
}

func NewMidtransUsecase(userRepo domain.UserRepository, topUpRepo domain.TopUpRepository, transactor helpers.Transactor, midTransServerKey string) domain.MidtransUsecase {
	return &midtransUsecase{
		userRepo:          userRepo,
		topUpRepo:         topUpRepo,
		transactor:        transactor,
		midtransServerKey: midTransServerKey,
	}
}

func (usecase *midtransUsecase) InitiateTopUp(ctx context.Context, userID string, req dto.TopUpRequestDTO) (*dto.TopUpResponseDTO, error) {
	if req.Nominal <= 0 {
		return nil, errors.New("nominal top up harus lebih besar dari 0")
	}

	orderID := fmt.Sprintf("TP-%d-%s", time.Now().Unix(), userID[:8])

	midtransPayload := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id":     orderID,
			"gross_amount": req.Nominal,
		},
	}

	payloadBytes, _ := json.Marshal(midtransPayload)

	snapURL := "https://app.sandbox.midtrans.com/snap/v1/transactions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", snapURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	serverKey := usecase.midtransServerKey
	authHeader := base64.StdEncoding.EncodeToString([]byte(serverKey + ":"))
	httpReq.Header.Set("Authorization", "Basic "+authHeader)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		logrus.WithError(err).Error("usecase.InitiateTopUp: failed to call midtrans API")
		return nil, err
	}
	defer resp.Body.Close()

	type MidtransSnapResponse struct {
		Token       string `json:"token"`
		RedirectURL string `json:"redirect_url"`
	}

	var snapResp MidtransSnapResponse
	if err := json.NewDecoder(resp.Body).Decode(&snapResp); err != nil {
		return nil, err
	}

	topUpData := &domain.TopUp{
		UserID:                userID,
		Nominal:               req.Nominal,
		MetodePembayaran:      req.MetodePembayaran,
		MidtransOrderID:       orderID,
		MidtransTransactionID: orderID,
		Status:                "pending",
		SnapToken:             snapResp.Token,
		SnapRedirectURL:       snapResp.RedirectURL,
	}

	if err := usecase.topUpRepo.CreateTopUp(ctx, topUpData); err != nil {
		logrus.WithError(err).Error("usecase.InitiateTopUp: failed to save top up record")
		return nil, err
	}

	response := &dto.TopUpResponseDTO{
		ID: topUpData.ID,
		User: &dto.UserResponseDTO{
			ID: userID,
		},
		Nominal:               topUpData.Nominal,
		MetodePembayaran:      topUpData.MetodePembayaran,
		MidtransOrderID:       topUpData.MidtransOrderID,
		MidtransTransactionID: topUpData.MidtransOrderID,
		Status:                topUpData.Status,
		SnapToken:             topUpData.SnapToken,
		SnapRedirectURL:       topUpData.SnapRedirectURL,
	}
	return response, nil
}

func (usecase *midtransUsecase) TopUpSaldo(ctx context.Context, cb dto.MidtransCallbackDTO) (bool, error) {
	sha := sha512.New()
	payloadSignature := cb.OrderID + cb.StatusCode + cb.GrossAmount + usecase.midtransServerKey
	sha.Write([]byte(payloadSignature))
	expectedSignature := hex.EncodeToString(sha.Sum(nil))

	if cb.SignatureKey != expectedSignature {
		logrus.Error("usecase.TopUpSaldo: alert! invalid midtrans signature key detected")
		return false, errors.New("invalid signature key")
	}
	err := usecase.transactor.WithTransaction(ctx, func(txCtx context.Context) error {

		topUp, err := usecase.topUpRepo.GetTopUpByMidtransOrderID(txCtx, cb.OrderID)
		if err != nil {
			return err
		}
		if topUp == nil {
			return errors.New("top-up record not found")
		}

		if topUp.Status == "settlement" || topUp.Status == "success" {
			logrus.Info("usecase.TopUpSaldo: transaction already processed previously")
			return nil
		}

		user := topUp.User
		if user == nil || user.ID == "" {
			return errors.New("user owner of transaction not found")
		}

		rawCallback, _ := json.Marshal(cb)
		topUp.RawCallBack = rawCallback
		topUp.MidtransTransactionID = cb.TransactionID
		topUp.Status = "sukses"

		switch cb.TransactionStatus {
		case "settlement", "capture":
			user.Saldo += topUp.Nominal
			user.UpdatedAt = time.Now()
			topUp.PaidAt = time.Now()

			err = usecase.userRepo.AddUserSaldo(txCtx, user.ID, topUp.Nominal)
			if err != nil {
				return err
			}

		case "pending", "deny", "expire", "cancel":
			logrus.Infof("usecase.TopUpSaldo: transaction status updated to %s", cb.TransactionStatus)
		}

		err = usecase.topUpRepo.UpdateTopUp(txCtx, topUp)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logrus.WithError(err).Error("usecase.TopUpSaldo: database transaction failed")
		return false, err
	}

	return true, nil
}

func (usecase *midtransUsecase) GetTopUpStatus(ctx context.Context, orderID string) (*dto.TopUpResponseDTO, error) {
	topUp, err := usecase.topUpRepo.GetTopUpByMidtransOrderID(ctx, orderID)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetTopUpStatus: failed to get top-up record")
		return nil, err
	}
	if topUp == nil {
		logrus.Info("usecase.GetTopUpStatus: top-up record not found for order ID: ", orderID)
		return nil, nil
	}

	response := &dto.TopUpResponseDTO{
		ID: topUp.ID,
		User: &dto.UserResponseDTO{
			ID:    topUp.User.ID,
			Name:  topUp.User.Name,
			Email: topUp.User.Email,
			Saldo: topUp.User.Saldo,
		},
		Nominal:               topUp.Nominal,
		MetodePembayaran:      topUp.MetodePembayaran,
		MidtransOrderID:       topUp.MidtransOrderID,
		MidtransTransactionID: topUp.MidtransTransactionID,
		Status:                topUp.Status,
		SnapToken:             topUp.SnapToken,
		SnapRedirectURL:       topUp.SnapRedirectURL,
		RawCallBack:           string(topUp.RawCallBack),
		PaidAt:                topUp.PaidAt,
	}

	return response, nil
}
