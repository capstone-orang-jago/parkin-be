package usecase

import (
	"bytes"
	"context"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
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
	statusProvider    MidtransStatusProvider
}

func NewMidtransUsecase(userRepo domain.UserRepository, topUpRepo domain.TopUpRepository, transactor helpers.Transactor, midTransServerKey string, statusProvider MidtransStatusProvider) domain.MidtransUsecase {
	return &midtransUsecase{
		userRepo:          userRepo,
		topUpRepo:         topUpRepo,
		transactor:        transactor,
		midtransServerKey: midTransServerKey,
		statusProvider:    statusProvider,
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
	if cb.OrderID == "" {
		return false, errors.New("order id is required")
	}

	sha := sha512.New()
	payloadSignature := cb.OrderID + cb.StatusCode + cb.GrossAmount + usecase.midtransServerKey
	sha.Write([]byte(payloadSignature))
	expectedSignature := hex.EncodeToString(sha.Sum(nil))

	if subtle.ConstantTimeCompare([]byte(cb.SignatureKey), []byte(expectedSignature)) != 1 {
		logrus.Error("usecase.TopUpSaldo: alert! invalid midtrans signature key detected")
		return false, errors.New("invalid signature key")
	}

	status, err := usecase.statusProvider.GetTransactionStatus(ctx, cb.OrderID)
	if err != nil {
		logrus.WithError(err).Error("usecase.TopUpSaldo: failed to verify transaction status with midtrans")
		return false, errors.New("failed to verify transaction status with midtrans")
	}
	if status == nil {
		return false, errors.New("midtrans transaction status not found")
	}
	if status.OrderID != cb.OrderID {
		return false, errors.New("order id does not match midtrans transaction")
	}
	if status.StatusCode != cb.StatusCode {
		return false, errors.New("status code does not match midtrans transaction")
	}
	if status.TransactionID != cb.TransactionID {
		return false, errors.New("transaction id does not match midtrans transaction")
	}

	newStatus, recognised := mapMidtransStatus(status.TransactionStatus, status.FraudStatus)
	if !recognised {
		return false, fmt.Errorf("unrecognised midtrans transaction status: %s", status.TransactionStatus)
	}

	err = usecase.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		topUp, err := usecase.topUpRepo.GetTopUpByMidtransOrderIDForUpdate(txCtx, cb.OrderID)
		if err != nil {
			return err
		}
		if topUp == nil {
			return errors.New("top-up record not found")
		}

		grossAmount, err := parseGrossAmount(status.GrossAmount)
		if err != nil {
			return errors.New("invalid gross amount from midtrans")
		}
		if grossAmount != topUp.Nominal {
			return errors.New("gross amount does not match top-up nominal")
		}
		if topUp.MidtransTransactionID != "" && topUp.MidtransTransactionID != topUp.MidtransOrderID && topUp.MidtransTransactionID != status.TransactionID {
			return errors.New("top-up transaction id does not match midtrans transaction")
		}

		if topUp.Status == "sukses" {
			logrus.Info("usecase.TopUpSaldo: transaction already processed previously")
			return nil
		}
		if topUp.Status != "pending" {
			logrus.WithFields(logrus.Fields{
				"order_id":        topUp.MidtransOrderID,
				"stored_status":   topUp.Status,
				"midtrans_status": status.TransactionStatus,
			}).Info("usecase.TopUpSaldo: ignoring final status change for already finalized top-up")
			return nil
		}

		rawCallback, _ := json.Marshal(cb)
		topUp.RawCallBack = rawCallback
		topUp.MidtransTransactionID = status.TransactionID
		topUp.Status = newStatus

		if newStatus == "sukses" {
			topUp.PaidAt = time.Now()

			user := topUp.User
			if user == nil || user.ID == "" {
				return errors.New("user owner of transaction not found")
			}

			if err := usecase.userRepo.AddUserSaldo(txCtx, user.ID, topUp.Nominal); err != nil {
				logrus.WithError(err).Error("usecase.TopUpSaldo: failed to add user saldo")
				return err
			}
		}

		if err := usecase.topUpRepo.UpdateTopUp(txCtx, topUp); err != nil {
			logrus.WithError(err).Error("usecase.TopUpSaldo: failed to update top-up record")
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

func (usecase *midtransUsecase) GetTopUpStatus(ctx context.Context, orderID, userID string) (*dto.TopUpResponseDTO, error) {
	topUp, err := usecase.topUpRepo.GetTopUpByMidtransOrderID(ctx, orderID)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetTopUpStatus: failed to get top-up record")
		return nil, err
	}
	if topUp == nil {
		logrus.Info("usecase.GetTopUpStatus: top-up record not found for order ID: ", orderID)
		return nil, nil
	}
	if topUp.User == nil || topUp.User.ID != userID {
		return nil, helpers.ErrForbidden
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
		PaidAt:                topUp.PaidAt,
	}

	return response, nil
}

func mapMidtransStatus(transactionStatus, fraudStatus string) (string, bool) {
	switch transactionStatus {
	case "settlement":
		return "sukses", true
	case "capture":
		if fraudStatus == "challenge" {
			return "pending", true
		}
		return "sukses", true
	case "expire":
		return "expired", true
	case "deny", "cancel":
		return "gagal", true
	case "pending":
		return "pending", true
	default:
		return "", false
	}
}

func parseGrossAmount(amount string) (int64, error) {
	value, ok := new(big.Rat).SetString(amount)
	if !ok {
		return 0, errors.New("invalid gross amount")
	}
	if !value.IsInt() {
		return 0, errors.New("gross amount must be a whole number")
	}
	if !value.Num().IsInt64() {
		return 0, errors.New("gross amount is out of range")
	}
	return value.Num().Int64(), nil
}
