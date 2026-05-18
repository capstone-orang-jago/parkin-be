package app

import (
	"github.com/capstone-orang-jago/parkin-be/config"
	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/capstone-orang-jago/parkin-be/http/handler"
	"github.com/capstone-orang-jago/parkin-be/repository"
	"github.com/capstone-orang-jago/parkin-be/usecase"
)

type App struct {
	UserHandler               domain.UserHandler
	UserUsecase               domain.UserUsecase
	UserRepository            domain.UserRepository
	AuthHandler               domain.AuthHandler
	AuthUsecase               domain.AuthUsecase
	AuthRepository            domain.AuthRepository
	LokasiParkirRepository    domain.LokasiParkirRepository
	LokasiParkirUsecase       domain.LokasiParkirUsecase
	LokasiParkirHandler       domain.LokasiParkirHandler
	QRSessionHandler          domain.QRSessionHandler
	QRSessionUsecase          domain.QRSessionUsecase
	QRSessionRepository       domain.QRSessionRepository
	TransaksiParkirRepository domain.TransaksiParkirRepository
	TransaksiParkirUsecase    domain.TransaksiParkirUsecase
	TransaksiParkirHandler    domain.TransaksiParkirHandler
	OtpCodeHandler            domain.OtpCodeHandler
	OtpCodeUsecase            domain.OtpCodeUsecase
	OtpCodeRepository         domain.OtpCodeRepository
	TopUpRepository           domain.TopUpRepository
	MidtransUsecase           domain.MidtransUsecase
	MidtransHandler           domain.MidtransHandler
	DashboardHandler          domain.DashboardHandler
	DashboardUsecase          domain.DashboardUsecase
	DashboardRepository       domain.DashboardRepository
}

func InitApp(cfg *config.Config, waService domain.WaService) App {

	db := config.DB
	transactor := helpers.NewTransactor(db)

	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	lokasiParkirRepo := repository.NewLokasiParkirRepository(db)
	qrSessionRepo := repository.NewQRSessionRepository(db)
	otpCodeRepo := repository.NewOtpCodeRepository(db)
	topUpRepo := repository.NewTopUpRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	otpCodeUsecase := usecase.NewOtpCodeUsecase(otpCodeRepo, userRepo, waService)
	authUsecase := usecase.NewAuthUsecase(authRepo, userRepo, transactor, otpCodeUsecase)
	lokasiParkirUsecase := usecase.NewLokasiParkirUsecase(lokasiParkirRepo)
	qrSessionUsecase := usecase.NewQRSessionUsecase(qrSessionRepo, lokasiParkirRepo, transactor, config.AppConfig.QR_SECRET_KEY)

	userHandler := handler.NewUserHandler(userUsecase)
	otpCodeHandler := handler.NewOtpCodeHandler(otpCodeUsecase)
	authHandler := handler.NewAuthHandler(authUsecase, userUsecase)
	lokasiParkirHandler := handler.NewLokasiParkirHandler(lokasiParkirUsecase)
	qrSessionHandler := handler.NewQRSessionHandler(qrSessionUsecase)

	transaksiParkirRepo := repository.NewTransaksiParkirRepository(db)
	transaksiParkirUsecase := usecase.NewTransaksiParkirUsecase(transaksiParkirRepo, userRepo, qrSessionRepo, lokasiParkirRepo, transactor)
	transaksiParkirHandler := handler.NewTransaksiParkirHandler(transaksiParkirUsecase)

	midtransUsecase := usecase.NewMidtransUsecase(userRepo, topUpRepo, transactor, config.AppConfig.MIDTRANS_SERVER_KEY)
	midtransHandler := handler.NewMidtransHandler(midtransUsecase)

	dashBoardRepo := repository.NewDashboardRepository(db)
	dashboardUsecase := usecase.NewDashboardUsecase(dashBoardRepo, transaksiParkirRepo, userRepo)
	dashboardHandler := handler.NewDashboardHandler(dashboardUsecase)
	return App{
		UserHandler:               userHandler,
		UserUsecase:               userUsecase,
		UserRepository:            userRepo,
		AuthHandler:               authHandler,
		AuthUsecase:               authUsecase,
		AuthRepository:            authRepo,
		OtpCodeHandler:            otpCodeHandler,
		OtpCodeUsecase:            otpCodeUsecase,
		OtpCodeRepository:         otpCodeRepo,
		LokasiParkirRepository:    lokasiParkirRepo,
		LokasiParkirUsecase:       lokasiParkirUsecase,
		LokasiParkirHandler:       lokasiParkirHandler,
		QRSessionHandler:          qrSessionHandler,
		QRSessionUsecase:          qrSessionUsecase,
		QRSessionRepository:       qrSessionRepo,
		TransaksiParkirRepository: transaksiParkirRepo,
		TransaksiParkirUsecase:    transaksiParkirUsecase,
		TransaksiParkirHandler:    transaksiParkirHandler,
		TopUpRepository:           topUpRepo,
		MidtransUsecase:           midtransUsecase,
		MidtransHandler:           midtransHandler,
		DashboardHandler:          dashboardHandler,
		DashboardUsecase:          dashboardUsecase,
		DashboardRepository:       dashBoardRepo,
	}
}
