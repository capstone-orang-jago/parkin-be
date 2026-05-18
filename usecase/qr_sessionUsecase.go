package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/sirupsen/logrus"
)

type qrSessionUsecase struct {
	qrSessionRepo    domain.QRSessionRepository
	lokasiParkirRepo domain.LokasiParkirRepository
	transactor       helpers.Transactor
	qrSecretKey      string
}

func NewQRSessionUsecase(qrSessionRepo domain.QRSessionRepository, lokasiParkirRepo domain.LokasiParkirRepository, transactor helpers.Transactor, qrSecretKey string) domain.QRSessionUsecase {
	return &qrSessionUsecase{
		qrSessionRepo:    qrSessionRepo,
		lokasiParkirRepo: lokasiParkirRepo,
		transactor:       transactor,
		qrSecretKey:      qrSecretKey,
	}
}

func (usecase *qrSessionUsecase) CreateQRSession(ctx context.Context, petugasID string) (res *dto.QRSessionResponseDTO, err error) {
	var result *domain.QRSession

	err = usecase.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		activeQR, err := usecase.qrSessionRepo.GetActiveQRByPetugasID(txCtx, petugasID)
		if err != nil {
			logrus.WithError(err).Error("usecase.CreateQRSession: failed to get active QR session by petugas ID")
			return err
		}

		if activeQR != nil {
			err = usecase.RevokeQRSession(txCtx, activeQR.ID)
			if err != nil {
				logrus.WithError(err).Error("usecase.CreateQRSession: failed to revoke existing QR session")
				return err
			}
		}

		lokasiParkir, err := usecase.lokasiParkirRepo.GetLokasiParkirByPetugasID(txCtx, petugasID)
		if err != nil {
			logrus.WithError(err).Error("usecase.CreateQRSession: failed to get lokasi by petugas ID")
			return err
		}
		if lokasiParkir == nil {
			return errors.New("usecase.CreateQRSession: no lokasi found for petugas ID")
		}

		now := time.Now()
		timestamp := now.Unix()
		message := fmt.Sprintf("parkin:%s:%s:%d", lokasiParkir.ID, petugasID, timestamp)

		hmac := hmac.New(sha256.New, []byte(usecase.qrSecretKey))
		hmac.Write([]byte(message))
		signature := hex.EncodeToString(hmac.Sum(nil))

		payload := fmt.Sprintf("%s:%s", message, signature)

		qrSession := &domain.QRSession{
			LokasiID:  lokasiParkir.ID,
			PetugasID: petugasID,
			Payload:   payload,
			HMAC:      signature,
			ExpiredAt: now.Add(5 * time.Minute),
			IsRevoked: false,
		}

		err = usecase.qrSessionRepo.CreateQRSession(txCtx, qrSession)
		if err != nil {
			logrus.WithError(err).Error("usecase.CreateQRSession: failed to create QR session")
			return err
		}

		result = qrSession
		return nil
	})
	if err != nil {
		logrus.WithError(err).Error("usecase.CreateQRSession: transaction failed")
		return nil, err
	}

	response := &dto.QRSessionResponseDTO{
		ID:        result.ID,
		Payload:   result.Payload,
		HMAC:      result.HMAC,
		ExpiredAt: result.ExpiredAt,
		IsRevoked: result.IsRevoked,
		LokasiParkir: &dto.LokasiParkirResponseDTO{
			ID:         result.LokasiParkir.ID,
			NamaLokasi: result.LokasiParkir.NamaLokasi,
			Alamat:     result.LokasiParkir.Alamat,
			IsActive:   result.LokasiParkir.IsActive,
			Petugas: &dto.UserResponseDTO{
				ID:    result.Petugas.ID,
				Name:  result.Petugas.Name,
				Email: result.Petugas.Email,
			},
		},
	}

	return response, nil
}

func (usecase *qrSessionUsecase) GetActiveQRByPetugasID(ctx context.Context, petugasID string) (res *dto.QRSessionResponseDTO, err error) {
	data, err := usecase.qrSessionRepo.GetActiveQRByPetugasID(ctx, petugasID)

	if err != nil {
		logrus.WithError(err).Error("usecase.GetActiveQRByPetugasID: failed to get active QR session by petugas ID")
		return nil, err
	}

	if data == nil {
		logrus.Info("usecase.GetActiveQRByPetugasID: no active QR session found for petugas ID")
		return nil, nil
	}

	response := &dto.QRSessionResponseDTO{
		ID:        data.ID,
		Payload:   data.Payload,
		HMAC:      data.HMAC,
		ExpiredAt: data.ExpiredAt,
		IsRevoked: data.IsRevoked,
		LokasiParkir: &dto.LokasiParkirResponseDTO{
			ID:         data.LokasiParkir.ID,
			NamaLokasi: data.LokasiParkir.NamaLokasi,
			Alamat:     data.LokasiParkir.Alamat,
			IsActive:   data.LokasiParkir.IsActive,
			Petugas: &dto.UserResponseDTO{
				ID:    data.Petugas.ID,
				Name:  data.Petugas.Name,
				Email: data.Petugas.Email,
			},
		},
	}

	return response, nil
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

func (usecase *qrSessionUsecase) GetQRSessionByPayload(ctx context.Context, payload string) (*dto.QRSessionResponseDTO, error) {
	if payload == "" {
		logrus.Error("usecase.GetQRSessionByPayload: payload is required")
		return nil, errors.New("payload is required")
	}

	qrSession, err := usecase.qrSessionRepo.GetQRSessionByPayload(ctx, payload)

	if err != nil {
		logrus.WithError(err).Error("usecase.GetQRSessionByPayload: failed to get QR session by payload")
		return nil, err
	}

	if qrSession == nil {
		logrus.Info("usecase.GetQRSessionByPayload: no QR session found for payload")
		return nil, nil
	}

	response := &dto.QRSessionResponseDTO{
		ID:        qrSession.ID,
		Payload:   qrSession.Payload,
		HMAC:      qrSession.HMAC,
		ExpiredAt: qrSession.ExpiredAt,
		IsRevoked: qrSession.IsRevoked,
		LokasiParkir: &dto.LokasiParkirResponseDTO{
			ID:         qrSession.LokasiParkir.ID,
			NamaLokasi: qrSession.LokasiParkir.NamaLokasi,
			Alamat:     qrSession.LokasiParkir.Alamat,
			IsActive:   qrSession.LokasiParkir.IsActive,
			Petugas: &dto.UserResponseDTO{
				ID:    qrSession.Petugas.ID,
				Name:  qrSession.Petugas.Name,
				Email: qrSession.Petugas.Email,
			},
		},
	}
	return response, nil
}
