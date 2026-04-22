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
}

func InitApp(cfg *config.Config) App {

	db := config.DB

	userRepo := repository.NewUserRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepo)

	userHandler := handler.NewUserHandler(userUsecase)
	return App{
		UserHandler:    userHandler,
		UserUsecase:    userUsecase,
		UserRepository: userRepo,
	}
}
