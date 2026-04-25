package app

import (
	"github.com/capstone-orang-jago/parkin-be/config"
	"github.com/capstone-orang-jago/parkin-be/domain"
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
}

func InitApp(cfg *config.Config) App {

	db := config.DB

	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	lokasiParkirRepo := repository.NewLokasiParkirRepository(db)
	qrSessionRepo := repository.NewQRSessionRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepo)
	authUsecase := usecase.NewAuthUsecase(authRepo, userRepo)
	lokasiParkirUsecase := usecase.NewLokasiParkirUsecase(lokasiParkirRepo)
	qrSessionUsecase := usecase.NewQRSessionUsecase(qrSessionRepo, lokasiParkirUsecase)

	userHandler := handler.NewUserHandler(userUsecase)
	authHandler := handler.NewAuthHandler(authUsecase, userUsecase)
	lokasiParkirHandler := handler.NewLokasiParkirHandler(lokasiParkirUsecase)
	qrSessionHandler := handler.NewQRSessionHandler(qrSessionUsecase)

	transaksiParkirRepo := repository.NewTransaksiParkirRepository(db)
	transaksiParkirUsecase := usecase.NewTransaksiParkirUsecase(transaksiParkirRepo, userUsecase, qrSessionUsecase, lokasiParkirUsecase)
	transaksiParkirHandler := handler.NewTransaksiParkirHandler(transaksiParkirUsecase)
	return App{
		UserHandler:               userHandler,
		UserUsecase:               userUsecase,
		UserRepository:            userRepo,
		AuthHandler:               authHandler,
		AuthUsecase:               authUsecase,
		AuthRepository:            authRepo,
		LokasiParkirRepository:    lokasiParkirRepo,
		LokasiParkirUsecase:       lokasiParkirUsecase,
		LokasiParkirHandler:       lokasiParkirHandler,
		QRSessionHandler:          qrSessionHandler,
		QRSessionUsecase:          qrSessionUsecase,
		QRSessionRepository:       qrSessionRepo,
		TransaksiParkirRepository: transaksiParkirRepo,
		TransaksiParkirUsecase:    transaksiParkirUsecase,
		TransaksiParkirHandler:    transaksiParkirHandler,
	}
}
