package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/sirupsen/logrus"
)

type otpCodeUsecase struct {
	otpCodeRepo domain.OtpCodeRepository
	userRepo    domain.UserRepository
	waService   domain.WaService
}

func NewOtpCodeUsecase(otpCodeRepo domain.OtpCodeRepository, userRepo domain.UserRepository, waService domain.WaService) domain.OtpCodeUsecase {
	return &otpCodeUsecase{otpCodeRepo: otpCodeRepo, userRepo: userRepo, waService: waService}
}

func (usecase *otpCodeUsecase) SendOtpCode(ctx context.Context, telephone string) error {
	telephone, err := helpers.NormalizeTelephone(telephone)
	if err != nil {
		return err
	}
	user, err := usecase.userRepo.GetUserByTelephone(ctx, telephone)
	if err != nil {
		logrus.WithError(err).Error("usecase.SendOtpCode: failed to get user by telephone")
		return err
	}
	if user == nil {
		logrus.WithField("telephone", telephone).Info("usecase.SendOtpCode: no user found for telephone")
		return errors.New("usecase.SendOtpCode: user not found")
	}

	lastOtp, err := usecase.otpCodeRepo.GetOtpCodeByLastActive(ctx, user.ID)
	if err != nil {
		logrus.WithError(err).Error("usecase.SendOtpCode: failed to get last active OTP code")
		return err
	}
	if lastOtp != nil {
		remainingTime := time.Since(lastOtp.CreatedAt)

		if remainingTime < 1*time.Minute {
			secondsLeft := int(60 - remainingTime.Seconds())
			logrus.WithField("userID", user.ID).Infof("usecase.SendOtpCode: OTP code already sent, please wait %d seconds before requesting a new one", secondsLeft)
			return errors.New("usecase.SendOtpCode: wait until the previous OTP code expires")
		}

		if !lastOtp.IsUsed {
			lastOtp.IsUsed = true
			usecase.otpCodeRepo.UpdateOtpCode(ctx, lastOtp)
		}
	}

	code := strconv.Itoa(rand.IntN(900000) + 100000)
	hashedCode := sha256.Sum256([]byte(code))

	otp := &domain.OtpCode{
		UserId:    user.ID,
		Code:      hex.EncodeToString(hashedCode[:]),
		ExpiredAt: time.Now().Add(5 * time.Minute),
		Attempts:  0,
		IsUsed:    false,
		CreatedAt: time.Now(),
	}

	if err := usecase.otpCodeRepo.CreateOtpCode(ctx, otp); err != nil {
		logrus.WithError(err).Error("usecase.SendOtpCode: failed to create OTP code")
		return err
	}

	err = usecase.waService.SendMessage(telephone, "Your OTP code is: "+code)
	if err != nil {
		logrus.WithError(err).Error("usecase.SendOtpCode: failed to send OTP code via WhatsApp")
		return err
	}

	return nil
}

func (usecase *otpCodeUsecase) VerifyOtpCode(ctx context.Context, telephone string, code string) (bool, error) {
	user, err := usecase.userRepo.GetUserByTelephone(ctx, telephone)
	if err != nil {
		logrus.WithError(err).Error("usecase.VerifyOtpCode: failed to get user by telephone")
		return false, err
	}
	if user == nil {
		logrus.WithField("telephone", telephone).Info("usecase.VerifyOtpCode: no user found for telephone")
		return false, nil
	}

	otp, err := usecase.otpCodeRepo.GetOtpCodeByLastActive(ctx, user.ID)
	if err != nil {
		logrus.WithError(err).Error("usecase.VerifyOtpCode: failed to get OTP code")
		return false, err
	}
	if otp == nil {
		logrus.WithField("userID", user.ID).Info("usecase.VerifyOtpCode: no OTP code found for user")
		return false, errors.New("usecase.VerifyOtpCOde: otp code not found")
	}
	if time.Now().After(otp.ExpiredAt) {
		logrus.WithField("userID", user.ID).Info("usecase.VerifyOtpCode: OTP code expired")
		usecase.otpCodeRepo.DeleteOtpCode(ctx, otp.ID)
		return false, errors.New("usecase.VerifyOtpCode: otp code expired")
	}

	hashedCode := sha256.Sum256([]byte(code))
	if otp.Code != hex.EncodeToString(hashedCode[:]) {
		otp.Attempts++
		if otp.Attempts >= 3 {
			err := usecase.otpCodeRepo.DeleteOtpCode(ctx, otp.ID)
			if err != nil {
				logrus.WithError(err).Error("usecase.VerifyOtpCode: failed to delete OTP code after attempts exceeded")
			}
			logrus.WithError(err).Error("usecase.VerifyOtpCode: OTP code attempts exceeded")
			return false, errors.New("invalid otp code")
		}
		if err := usecase.otpCodeRepo.UpdateOtpCode(ctx, otp); err != nil {
			logrus.WithError(err).Error("usecase.VerifyOtpCode: failed to update OTP code attempts")
		}

		logrus.WithField("userID", user.ID).Info("usecase.VerifyOtpCode: invalid OTP code")
		return false, errors.New("invalid OTP code")
	}

	otp.IsUsed = true
	usecase.otpCodeRepo.UpdateOtpCode(ctx, otp)

	return true, nil
}
