package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/capstone-orang-jago/parkin-be/config"
	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/sirupsen/logrus"
)

const qrReuseCooldown = 5 * time.Second

type transaksiParkirUsecase struct {
	transaksiParkirRepo domain.TransaksiParkirRepository
	userRepo            domain.UserRepository
	qrSessionRepo       domain.QRSessionRepository
	lokasiParkirRepo    domain.LokasiParkirRepository
	transactor          helpers.Transactor
}

func NewTransaksiParkirUsecase(transaksiParkirRepo domain.TransaksiParkirRepository, userRepo domain.UserRepository, qrSessionRepo domain.QRSessionRepository, lokasiParkirRepo domain.LokasiParkirRepository, transactor helpers.Transactor) domain.TransaksiParkirUsecase {
	return &transaksiParkirUsecase{
		transaksiParkirRepo: transaksiParkirRepo,
		userRepo:            userRepo,
		qrSessionRepo:       qrSessionRepo,
		lokasiParkirRepo:    lokasiParkirRepo,
		transactor:          transactor,
	}
}

func (usecase *transaksiParkirUsecase) CreateTransaksiParkir(ctx context.Context, req *dto.TransaksiParkirRequestDTO, userID string) (*dto.TransaksiParkirResponseDTO, error) {
	if req.QRPayload == "" {
		return nil, fmt.Errorf("%w: QR payload is required", helpers.ErrInvalidInput)
	}
	if req.JenisKendaraan == "" {
		return nil, fmt.Errorf("%w: jenis kendaraan is required", helpers.ErrInvalidInput)
	}

	parts := strings.Split(req.QRPayload, ":")
	if len(parts) != 5 {
		return nil, fmt.Errorf("%w: invalid QR payload format", helpers.ErrInvalidInput)
	}

	prefix := parts[0]
	lokasiID := parts[1]
	petugasID := parts[2]
	timestampStr := parts[3]
	signature := parts[4]

	if prefix != "parkin" {
		return nil, fmt.Errorf("%w: invalid QR payload prefix", helpers.ErrInvalidInput)
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid timestamp in QR payload", helpers.ErrInvalidInput)
	}

	const expirySeconds = 300
	if time.Now().Unix()-timestamp > expirySeconds {
		return nil, fmt.Errorf("%w: QR payload has expired", helpers.ErrInvalidInput)
	}

	messageToVerify := fmt.Sprintf("parkin:%s:%s:%s", lokasiID, petugasID, timestampStr)

	h := hmac.New(sha256.New, []byte(config.AppConfig.QR_SECRET_KEY))
	h.Write([]byte(messageToVerify))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return nil, fmt.Errorf("%w: invalid QR payload signature", helpers.ErrInvalidInput)
	}

	var tarif int64
	switch req.JenisKendaraan {
	case "motor":
		tarif = 2000
	case "mobil":
		tarif = 4000
	default:
		return nil, fmt.Errorf("%w: invalid jenis kendaraan", helpers.ErrInvalidInput)
	}

	var transaksiParkir *domain.TransaksiParkir

	err = usecase.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		qrSession, err := usecase.qrSessionRepo.GetQRSessionByPayloadForUpdate(txCtx, req.QRPayload)
		if err != nil {
			logrus.WithError(err).Error("usecase.CreateTransaksiParkir: failed to get QR session by payload")
			return err
		}
		if qrSession == nil {
			return fmt.Errorf("%w: invalid QR payload", helpers.ErrInvalidInput)
		}
		if qrSession.LokasiID != lokasiID || qrSession.PetugasID != petugasID {
			return fmt.Errorf("%w: QR payload does not match the active session", helpers.ErrInvalidInput)
		}
		if qrSession.Petugas == nil || qrSession.Petugas.ID != petugasID {
			return fmt.Errorf("%w: invalid petugas ID in QR payload", helpers.ErrInvalidInput)
		}
		if qrSession.LokasiParkir == nil || !qrSession.LokasiParkir.IsActive {
			return fmt.Errorf("%w: lokasi parkir is not active", helpers.ErrInvalidInput)
		}
		if !qrSession.LastUsedAt.IsZero() && time.Since(qrSession.LastUsedAt) < qrReuseCooldown {
			return fmt.Errorf("%w: QR code is being used, please try again", helpers.ErrConflict)
		}

		user, err := usecase.userRepo.GetUserByIDForUpdate(txCtx, userID)
		if err != nil {
			logrus.WithError(err).Error("usecase.CreateTransaksiParkir: failed to get user by ID")
			return err
		}
		if user == nil {
			return fmt.Errorf("%w: user not found", helpers.ErrNotFound)
		}

		lokasi, err := usecase.lokasiParkirRepo.GetLokasiParkirByIDForUpdate(txCtx, lokasiID)
		if err != nil {
			logrus.WithError(err).Error("usecase.CreateTransaksiParkir: failed to get lokasi by ID")
			return err
		}
		if lokasi == nil {
			return fmt.Errorf("%w: lokasi parkir not found", helpers.ErrInvalidInput)
		}
		if !lokasi.IsActive {
			return fmt.Errorf("%w: lokasi parkir is not active", helpers.ErrInvalidInput)
		}
		if lokasi.Petugas == nil || lokasi.Petugas.ID != petugasID || lokasi.Petugas.Role != "petugas" {
			return fmt.Errorf("%w: petugas is no longer assigned to this location", helpers.ErrInvalidInput)
		}

		if user.Saldo < tarif {
			return fmt.Errorf("%w: user has insufficient saldo", helpers.ErrInsufficientBalance)
		}

		transaksiParkir = &domain.TransaksiParkir{
			KodeTransaksi:  fmt.Sprintf("TRX-%d", time.Now().UnixNano()),
			PenggunaID:     user.ID,
			PetugasID:      petugasID,
			LokasiID:       lokasi.ID,
			QRSessionID:    qrSession.ID,
			JenisKendaraan: req.JenisKendaraan,
			Tarif:          tarif,
			SaldoSebelum:   user.Saldo,
			SaldoSesudah:   user.Saldo - tarif,
			Status:         "sukses",
		}

		if err = usecase.transaksiParkirRepo.CreateTransaksiParkir(txCtx, transaksiParkir); err != nil {
			logrus.WithError(err).Error("usecase.CreateTransaksiParkir: failed to create transaksi parkir")
			return err
		}

		if err = usecase.userRepo.DeductUserSaldo(txCtx, user.ID, tarif); err != nil {
			logrus.WithError(err).Error("usecase.CreateTransaksiParkir: failed to deduct user saldo")
			return err
		}

		if err = usecase.userRepo.AddUserSaldo(txCtx, petugasID, tarif); err != nil {
			logrus.WithError(err).Error("usecase.CreateTransaksiParkir: failed to add saldo to petugas")
			return err
		}

		if err = usecase.qrSessionRepo.MarkQRSessionUsed(txCtx, qrSession.ID); err != nil {
			logrus.WithError(err).Error("usecase.CreateTransaksiParkir: failed to mark QR session as used")
			return err
		}

		return nil
	})
	if err != nil {
		logrus.WithError(err).Error("usecase.CreateTransaksiParkir: transaction failed")
		return nil, err
	}

	res := &dto.TransaksiParkirResponseDTO{
		ID:            transaksiParkir.ID,
		KodeTransaksi: transaksiParkir.KodeTransaksi,
		Pengguna: &dto.UserResponseDTO{
			ID:    transaksiParkir.Pengguna.ID,
			Name:  transaksiParkir.Pengguna.Name,
			Email: transaksiParkir.Pengguna.Email,
		},
		Lokasi: &dto.LokasiParkirResponseDTO{
			ID:         transaksiParkir.LokasiParkir.ID,
			NamaLokasi: transaksiParkir.LokasiParkir.NamaLokasi,
			Alamat:     transaksiParkir.LokasiParkir.Alamat,
			IsActive:   transaksiParkir.LokasiParkir.IsActive,
			Petugas: &dto.UserResponseDTO{
				ID:    transaksiParkir.LokasiParkir.Petugas.ID,
				Name:  transaksiParkir.LokasiParkir.Petugas.Name,
				Email: transaksiParkir.LokasiParkir.Petugas.Email,
			},
		},
		JenisKendaraan: transaksiParkir.JenisKendaraan,
		Tarif:          transaksiParkir.Tarif,
		SaldoSebelum:   transaksiParkir.SaldoSebelum,
		SaldoSesudah:   transaksiParkir.SaldoSesudah,
		Status:         transaksiParkir.Status,
		CreatedAt:      helpers.FormatTime(transaksiParkir.CreatedAt),
	}

	return res, nil
}

func (usecase *transaksiParkirUsecase) GetTransaksiParkirByID(ctx context.Context, id, userID, role string) (*dto.TransaksiParkirResponseDTO, error) {
	transaksiParkir, err := usecase.transaksiParkirRepo.GetTransaksiParkirByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetTransaksiParkirByID: failed to get transaksi parkir by ID")
		return nil, err
	} else if transaksiParkir == nil {
		logrus.WithError(err).Info("usecase.GetTransaksiParkirByID: no transaksi parkir found for ID")
		return nil, nil
	}

	switch role {
	case "pengguna":
		if transaksiParkir.PenggunaID != userID {
			return nil, helpers.ErrForbidden
		}
	case "petugas":
		if transaksiParkir.PetugasID != userID {
			return nil, helpers.ErrForbidden
		}
	}

	res := &dto.TransaksiParkirResponseDTO{
		ID:            transaksiParkir.ID,
		KodeTransaksi: transaksiParkir.KodeTransaksi,
		Pengguna: &dto.UserResponseDTO{
			ID:    transaksiParkir.Pengguna.ID,
			Name:  transaksiParkir.Pengguna.Name,
			Email: transaksiParkir.Pengguna.Email,
		},
		Lokasi: &dto.LokasiParkirResponseDTO{
			ID:         transaksiParkir.LokasiParkir.ID,
			NamaLokasi: transaksiParkir.LokasiParkir.NamaLokasi,
			Alamat:     transaksiParkir.LokasiParkir.Alamat,
			IsActive:   transaksiParkir.LokasiParkir.IsActive,
			Petugas: &dto.UserResponseDTO{
				ID:    transaksiParkir.LokasiParkir.Petugas.ID,
				Name:  transaksiParkir.LokasiParkir.Petugas.Name,
				Email: transaksiParkir.LokasiParkir.Petugas.Email,
			},
		},
		JenisKendaraan: transaksiParkir.JenisKendaraan,
		Tarif:          transaksiParkir.Tarif,
		SaldoSebelum:   transaksiParkir.SaldoSebelum,
		SaldoSesudah:   transaksiParkir.SaldoSesudah,
		Status:         transaksiParkir.Status,
		CreatedAt:      helpers.FormatTime(transaksiParkir.CreatedAt),
	}

	return res, nil
}

func (usecase *transaksiParkirUsecase) GetAllTransaksiParkirByUserID(ctx context.Context, userID, role string, page, limit int, filter string) ([]*dto.TransaksiParkirResponseDTO, error) {
	transaksis, err := usecase.transaksiParkirRepo.GetAllTransaksiParkirByUserID(ctx, userID, role, page, limit, filter)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetAllTransaksiParkirByUserID: failed to get transaksi parkir by user ID")
		return nil, err
	} else if transaksis == nil {
		logrus.WithError(err).Info("usecase.GetAllTransaksiParkirByUserID: no transaksi parkir found for user ID")
		return nil, nil
	}

	var res []*dto.TransaksiParkirResponseDTO
	for _, transaksi := range transaksis {
		res = append(res, &dto.TransaksiParkirResponseDTO{
			ID:            transaksi.ID,
			KodeTransaksi: transaksi.KodeTransaksi,
			Pengguna: &dto.UserResponseDTO{
				ID:    transaksi.Pengguna.ID,
				Name:  transaksi.Pengguna.Name,
				Email: transaksi.Pengguna.Email,
			},
			Lokasi: &dto.LokasiParkirResponseDTO{
				ID:         transaksi.LokasiParkir.ID,
				NamaLokasi: transaksi.LokasiParkir.NamaLokasi,
				Alamat:     transaksi.LokasiParkir.Alamat,
				IsActive:   transaksi.LokasiParkir.IsActive,
				Petugas: &dto.UserResponseDTO{
					ID:    transaksi.LokasiParkir.Petugas.ID,
					Name:  transaksi.LokasiParkir.Petugas.Name,
					Email: transaksi.LokasiParkir.Petugas.Email,
				},
			},
			JenisKendaraan: transaksi.JenisKendaraan,
			Tarif:          transaksi.Tarif,
			SaldoSebelum:   transaksi.SaldoSebelum,
			SaldoSesudah:   transaksi.SaldoSesudah,
			Status:         transaksi.Status,
			CreatedAt:      helpers.FormatTime(transaksi.CreatedAt),
		})
	}

	return res, nil
}
