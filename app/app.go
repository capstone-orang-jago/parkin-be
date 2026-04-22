package app

import (
	"github.com/capstone-orang-jago/parkin-be/config"
	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/http/handler"
	"github.com/capstone-orang-jago/parkin-be/repository"
	"github.com/capstone-orang-jago/parkin-be/usecase"
)

type App struct {
	UserHandler    domain.UserHandler
	UserUsecase    domain.UserUsecase
	UserRepository domain.UserRepository
	AuthHandler    domain.AuthHandler
	AuthUsecase    domain.AuthUsecase
	AuthRepository domain.AuthRepository
}

func InitApp(cfg *config.Config) App {

	db := config.DB

	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepo)
	authUsecase := usecase.NewAuthUsecase(authRepo, userRepo)

	userHandler := handler.NewUserHandler(userUsecase)
	authHandler := handler.NewAuthHandler(authUsecase, userUsecase)
	return App{
		UserHandler:    userHandler,
		UserUsecase:    userUsecase,
		UserRepository: userRepo,
		AuthHandler:    authHandler,
		AuthUsecase:    authUsecase,
		AuthRepository: authRepo,
	}
}
