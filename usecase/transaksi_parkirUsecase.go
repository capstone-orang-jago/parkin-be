package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/capstone-orang-jago/parkin-be/config"
	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/sirupsen/logrus"
)

type transaksiParkirUsecase struct {
	transaksiParkirRepo domain.TransaksiParkirRepository
	userUsecase         domain.UserUsecase
	qrSessionUsecase    domain.QRSessionUsecase
	lokasiParkirUsecase domain.LokasiParkirUsecase
}

func NewTransaksiParkirUsecase(transaksiParkirRepo domain.TransaksiParkirRepository, userUsecase domain.UserUsecase, qrSessionUsecase domain.QRSessionUsecase, lokasiParkirUsecase domain.LokasiParkirUsecase) domain.TransaksiParkirUsecase {
	return &transaksiParkirUsecase{
		transaksiParkirRepo: transaksiParkirRepo,
		userUsecase:         userUsecase,
		qrSessionUsecase:    qrSessionUsecase,
		lokasiParkirUsecase: lokasiParkirUsecase,
	}
}

func (usecase *transaksiParkirUsecase) CreateTransaksiParkir(ctx context.Context, req *dto.TransaksiParkirRequestDTO, userID string) (*dto.TransaksiParkirResponseDTO, error) {
	user, err := usecase.userUsecase.GetUserByID(ctx, userID)
	if err != nil {
		logrus.WithError(err).Error("usecase.CreateTransaksiParkir: failed to get user by ID")
		return nil, err
	} else if user == nil {
		logrus.WithError(err).Info("usecase.CreateTransaksiParkir: no user found for ID")
		return nil, nil
	}

	if req.QRPayload == "" {
		return nil, errors.New("QR payload is required")
	}
	if req.JenisKendaraan == "" {
		return nil, errors.New("jenis kendaraan is required")
	}

	parts := strings.Split(req.QRPayload, ":")
	if len(parts) != 5 {
		return nil, errors.New("invalid QR payload format")
	}

	prefix := parts[0]
	lokasiID := parts[1]
	petugasID := parts[2]
	timestampStr := parts[3]
	signature := parts[4]

	if prefix != "parkin" {
		logrus.WithError(err).Error("usecase.CreateTransaksiParkir: invalid QR payload prefix")
		return nil, errors.New("invalid QR payload prefix")
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		logrus.WithError(err).Error("usecase.CreateTransaksiParkir: invalid timestamp in QR payload")
		return nil, errors.New("invalid timestamp in QR payload")
	}

	const expirySeconds = 300
	if time.Now().Unix()-timestamp > expirySeconds {
		logrus.WithError(err).Error("usecase.CreateTransaksiParkir: QR payload has expired")
		return nil, errors.New("QR payload has expired")
	}

	messageToVerify := fmt.Sprintf("parkin:%s:%s:%s", lokasiID, petugasID, timestampStr)

	h := hmac.New(sha256.New, []byte(config.AppConfig.QR_SECRET_KEY))
	h.Write([]byte(messageToVerify))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		logrus.WithError(err).Error("usecase.CreateTransaksiParkir: invalid QR payload signature")
		return nil, errors.New("invalid QR payload signature")
	}

	qrSession, err := usecase.qrSessionUsecase.GetQRSessionByPayload(ctx, req.QRPayload)
	if err != nil {
		logrus.WithError(err).Error("usecase.CreateTransaksiParkir: failed to get QR session by payload")
		return nil, err
	}

	if qrSession == nil {
		return nil, errors.New("invalid QR payload")
	}

	var tarif int64
	switch req.JenisKendaraan {
	case "motor":
		tarif = 2000
	case "mobil":
		tarif = 5000
	default:
		return nil, errors.New("invalid jenis kendaraan")
	}

	petugas, err := usecase.userUsecase.GetUserByID(ctx, petugasID)
	if err != nil {
		logrus.WithError(err).Error("usecase.CreateTransaksiParkir: failed to get petugas by ID")
		return nil, err
	}

	if petugas == nil || petugas.Role != "petugas" {
		logrus.WithError(err).Error("usecase.CreateTransaksiParkir: invalid petugas ID in QR payload")
		return nil, errors.New("invalid petugas ID in QR payload")
	}

	lokasi, err := usecase.lokasiParkirUsecase.GetLokasiParkirByID(ctx, lokasiID)
	if err != nil {
		logrus.WithError(err).Error("usecase.CreateTransaksiParkir: failed to get lokasi by ID")
		return nil, err
	}

	if user.Saldo < tarif {
		logrus.WithError(err).Error("usecase.CreateTransaksiParkir: user has insufficient saldo")
		return nil, errors.New("user has insufficient saldo")
	}

	transaksiParkir := &domain.TransaksiParkir{
		KodeTransaksi:  fmt.Sprintf("TRX-%d", time.Now().UnixNano()),
		PenggunaID:     user.ID,
		PetugasID:      petugas.ID,
		LokasiID:       lokasiID,
		JenisKendaraan: req.JenisKendaraan,
		Tarif:          tarif,
		SaldoSebelum:   user.Saldo,
		SaldoSesudah:   user.Saldo - tarif,
		Status:         "sukses",
	}

	err = usecase.transaksiParkirRepo.CreateTransaksiParkir(ctx, transaksiParkir)
	if err != nil {
		logrus.WithError(err).Error("usecase.CreateTransaksiParkir: failed to create transaksi parkir")
		return nil, err
	}

	err = usecase.userUsecase.DeductUserSaldo(ctx, user.ID, tarif)
	if err != nil {
		logrus.WithError(err).Error("usecase.CreateTransaksiParkir: failed to update user saldo")
		return nil, err
	}

	err = usecase.qrSessionUsecase.RevokeQRSession(ctx, qrSession.ID)
	if err != nil {
		logrus.WithError(err).Error("usecase.CreateTransaksiParkir: failed to revoke QR session")
		return nil, err
	}

	res := &dto.TransaksiParkirResponseDTO{
		ID:             transaksiParkir.ID,
		KodeTransaksi:  transaksiParkir.KodeTransaksi,
		Petugas:        petugas,
		Lokasi:         lokasi,
		JenisKendaraan: transaksiParkir.JenisKendaraan,
		Tarif:          transaksiParkir.Tarif,
		SaldoSebelum:   transaksiParkir.SaldoSebelum,
		SaldoSesudah:   transaksiParkir.SaldoSesudah,
		Status:         transaksiParkir.Status,
	}

	return res, nil
}

func (usecase *transaksiParkirUsecase) GetTransaksiParkirByID(ctx context.Context, id string) (*dto.TransaksiParkirResponseDTO, error) {
	transaksiParkir, err := usecase.transaksiParkirRepo.GetTransaksiParkirByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetTransaksiParkirByID: failed to get transaksi parkir by ID")
		return nil, err
	} else if transaksiParkir == nil {
		logrus.WithError(err).Info("usecase.GetTransaksiParkirByID: no transaksi parkir found for ID")
		return nil, nil
	}

	petugas, err := usecase.userUsecase.GetUserByID(ctx, transaksiParkir.PetugasID)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetTransaksiParkirByID: failed to get petugas by ID")
		return nil, err
	} else if petugas == nil {
		logrus.WithError(err).Info("usecase.GetTransaksiParkirByID: no petugas found for ID")
		return nil, nil
	}

	lokasi, err := usecase.lokasiParkirUsecase.GetLokasiParkirByID(ctx, transaksiParkir.LokasiID)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetTransaksiParkirByID: failed to get lokasi by ID")
		return nil, err
	} else if lokasi == nil {
		logrus.WithError(err).Info("usecase.GetTransaksiParkirByID: no lokasi found for ID")
		return nil, nil
	}

	res := &dto.TransaksiParkirResponseDTO{
		ID:             transaksiParkir.ID,
		KodeTransaksi:  transaksiParkir.KodeTransaksi,
		Petugas:        petugas,
		Lokasi:         lokasi,
		JenisKendaraan: transaksiParkir.JenisKendaraan,
		Tarif:          transaksiParkir.Tarif,
		SaldoSebelum:   transaksiParkir.SaldoSebelum,
		SaldoSesudah:   transaksiParkir.SaldoSesudah,
		Status:         transaksiParkir.Status,
	}

	return res, nil
}

func (usecase *transaksiParkirUsecase) GetAllTransaksiParkirByPenggunaID(ctx context.Context, penggunaID string) ([]*dto.TransaksiParkirResponseDTO, error) {
	transaksis, err := usecase.transaksiParkirRepo.GetAllTransaksiParkirByPenggunaID(ctx, penggunaID)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetAllTransaksiParkirByPenggunaID: failed to get transaksi parkir by pengguna ID")
		return nil, err
	} else if transaksis == nil {
		logrus.WithError(err).Info("usecase.GetAllTransaksiParkirByPenggunaID: no transaksi parkir found for pengguna ID")
		return nil, nil
	}

	var res []*dto.TransaksiParkirResponseDTO
	for _, transaksi := range transaksis {
		petugas, err := usecase.userUsecase.GetUserByID(ctx, transaksi.PetugasID)
		if err != nil {
			logrus.WithError(err).Error("usecase.GetAllTransaksiParkirByPenggunaID: failed to get petugas by ID")
			return nil, err
		} else if petugas == nil {
			logrus.WithError(err).Info("usecase.GetAllTransaksiParkirByPenggunaID: no petugas found for ID")
			return nil, nil
		}

		lokasi, err := usecase.lokasiParkirUsecase.GetLokasiParkirByID(ctx, transaksi.LokasiID)
		if err != nil {
			logrus.WithError(err).Error("usecase.GetAllTransaksiParkirByPenggunaID: failed to get lokasi by ID")
			return nil, err
		} else if lokasi == nil {
			logrus.WithError(err).Info("usecase.GetAllTransaksiParkirByPenggunaID: no lokasi found for ID")
			return nil, nil
		}

		res = append(res, &dto.TransaksiParkirResponseDTO{
			ID:             transaksi.ID,
			KodeTransaksi:  transaksi.KodeTransaksi,
			Petugas:        petugas,
			Lokasi:         lokasi,
			JenisKendaraan: transaksi.JenisKendaraan,
			Tarif:          transaksi.Tarif,
			SaldoSebelum:   transaksi.SaldoSebelum,
			SaldoSesudah:   transaksi.SaldoSesudah,
			Status:         transaksi.Status,
		})
	}

	return res, nil
}
