package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/capstone-orang-jago/parkin-be/config"
	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/sirupsen/logrus"
)

type qrSessionUsecase struct {
	qrSessionRepo       domain.QRSessionRepository
	parkirLokasiUsecase domain.LokasiParkirUsecase
}

func NewQRSessionUsecase(qrSessionRepo domain.QRSessionRepository, parkirLokasi domain.LokasiParkirUsecase) domain.QRSessionUsecase {
	return &qrSessionUsecase{
		qrSessionRepo:       qrSessionRepo,
		parkirLokasiUsecase: parkirLokasi,
	}
}

func (usecase *qrSessionUsecase) CreateQRSession(ctx context.Context, petugasID string) (res *dto.QRSessionResponseDTO, err error) {
	activeQR, err := usecase.GetActiveQRByPetugasID(ctx, petugasID)
	if err != nil {
		logrus.WithError(err).Error("usecase.CreateQRSession: failed to get active QR session by petugas ID")
		return nil, err
	}

	if activeQR != nil {
		err = usecase.RevokeQRSession(ctx, activeQR.ID)
		if err != nil {
			logrus.WithError(err).Error("usecase.CreateQRSession: failed to revoke existing QR session")
			return nil, err
		}
	}

	lokasiParkir, err := usecase.parkirLokasiUsecase.GetLokasiParkirByPetugasID(ctx, petugasID)
	if err != nil {
		logrus.WithError(err).Error("usecase.CreateQRSession: failed to get lokasi by petugas ID")
		return nil, err
	}

	timestamp := time.Now().Unix()
	message := fmt.Sprintf("parkin:%s:%s:%d", lokasiParkir.ID, petugasID, timestamp)

	hmac := hmac.New(sha256.New, []byte(config.AppConfig.QR_SECRET_KEY))
	hmac.Write([]byte(message))
	signature := hex.EncodeToString(hmac.Sum(nil))

	payload := fmt.Sprintf("%s:%s", message, signature)

	qrSession := &domain.QRSession{
		LocationID: lokasiParkir.ID,
		PetugasID:  petugasID,
		Payload:    payload,
		HMAC:       signature,
		ExpiresAt:  time.Now().Add(5 * time.Minute),
		IsRevoked:  false,
	}

	err = usecase.qrSessionRepo.CreateQRSession(ctx, qrSession)
	if err != nil {
		logrus.WithError(err).Error("usecase.CreateQRSession: failed to create QR session")
		return nil, err
	}

	response := &dto.QRSessionResponseDTO{
		ID:         qrSession.ID,
		LocationID: qrSession.LocationID,
		PetugasID:  qrSession.PetugasID,
		Payload:    qrSession.Payload,
		HMAC:       qrSession.HMAC,
		ExpiresAt:  qrSession.ExpiresAt,
		IsRevoked:  qrSession.IsRevoked,
	}

	return response, nil
}

func (usecase *qrSessionUsecase) GetActiveQRByPetugasID(ctx context.Context, petugasID string) (res *dto.QRSessionResponseDTO, err error) {
	data, err := usecase.qrSessionRepo.GetActiveQRByPetugasID(ctx, petugasID)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetActiveQRByPetugasID: failed to get active QR session by petugas ID")
		return nil, err
	}

	response := &dto.QRSessionResponseDTO{
		ID:         data.ID,
		LocationID: data.LocationID,
		PetugasID:  data.PetugasID,
		Payload:    data.Payload,
		HMAC:       data.HMAC,
		ExpiresAt:  data.ExpiresAt,
		IsRevoked:  data.IsRevoked,
	}

	return response, nil
}

func (usecase *qrSessionUsecase) GetLokasiByPetugasID(ctx context.Context, petugasID string) (string, error) {
	data, err := usecase.qrSessionRepo.GetLokasiByPetugasID(ctx, petugasID)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetLokasiByPetugasID: failed to get lokasi by petugas ID")
		return "", err
	}

	return data, nil
}

func (usecase *qrSessionUsecase) RevokeQRSession(ctx context.Context, id string) error {
	if id == "" {
		logrus.Error("usecase.RevokeQRSession: QR session ID is required")
		return errors.New("QR session id is required")
	}

	err := usecase.qrSessionRepo.RevokeQRSession(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("usecase.RevokeQRSession: failed to revoke QR session")
		return err
	}

	return nil
}
