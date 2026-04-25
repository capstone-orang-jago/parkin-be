package usecase

import (
	"context"
	"errors"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/sirupsen/logrus"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (usecase *userUsecase) GetAllUsers(ctx context.Context) (users []*domain.User, err error) {
	users, err = usecase.userRepo.GetAllUsers(ctx)

	if err != nil {
		logrus.WithError(err).Error("usecase.GetAllUsers: failed to get all users")
		return nil, err
	}

	return users, nil
}

func (usecase *userUsecase) GetUserByID(ctx context.Context, id string) (user *domain.User, err error) {
	user, err = usecase.userRepo.GetUserByID(ctx, id)

	if err != nil {
		logrus.WithError(err).Error("usecase.GetUserByID: failed to get user by ID")
		return nil, err
	}

	return user, nil
}

func (usecase *userUsecase) GetUserByEmail(ctx context.Context, email string) (user *domain.User, err error) {
	user, err = usecase.userRepo.GetUserByEmail(ctx, email)

	if err != nil {
		logrus.WithError(err).Error("usecase.GetUserByEmail: failed to get user by email")
		return nil, err
	}

	return user, nil
}

func (usecase *userUsecase) CreateUser(ctx context.Context, req dto.UserRequestDTO) (res *domain.User, err error) {
	if req.Name == "" {
		return nil, errors.New("name are required")
	}
	if req.Telephone == "" {
		return nil, errors.New("telephone are required")
	}
	if req.Email == "" {
		return nil, errors.New("email are required")
	}
	if req.Password == "" {
		return nil, errors.New("password are required")
	}

	if req.Password != "" {
		hashedPassword, err := helpers.HassPass(req.Password)
		if err != nil {
			logrus.Error("usecase.CreateUser: error hash password: ", err)
			return nil, err
		}
		req.Password = hashedPassword
	}

	user := &domain.User{
		Name:       req.Name,
		Telephone:  req.Telephone,
		Email:      req.Email,
		Password:   req.Password,
		Role:       "pengguna",
		Saldo:      0,
		IsVerified: false,
	}

	err = usecase.userRepo.CreateUser(ctx, user)
	if err != nil {
		logrus.WithError(err).Error("usecase.CreateUser: failed to create user")
		return nil, err
	}

	return user, nil
}

func (usecase *userUsecase) UpdateUser(ctx context.Context, req dto.UserRequestDTO, id string) (res *domain.User, err error) {
	user, err := usecase.userRepo.GetUserByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("usecase.UpdateUser: failed to get user by ID")
		return nil, err
	}

	if req.Name == "" && req.Email == "" && req.Password == "" && req.Telephone == "" {
		return nil, errors.New("at least one field must be provided to update")
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Telephone != "" {
		user.Telephone = req.Telephone
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := helpers.HassPass(req.Password)
		if err != nil {
			logrus.Error("usecase.UpdateUser: error hash password: ", err)
			return nil, err
		}
		user.Password = hashedPassword
	}

	err = usecase.userRepo.UpdateUser(ctx, user, id)
	if err != nil {
		logrus.WithError(err).Error("usecase.UpdateUser: failed to update user")
		return nil, err
	}

	return user, nil
}

func (usecase *userUsecase) DeleteUser(ctx context.Context, id string) error {
	_, err := usecase.userRepo.GetUserByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("usecase.DeleteUser: failed to get user by ID")
		return err
	}

	err = usecase.userRepo.DeleteUser(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("usecase.DeleteUser: failed to delete user")
		return err
	}
	return nil
}

func (usecase *userUsecase) GetUserSaldo(ctx context.Context, id string) (saldo int64, err error) {
	user, err := usecase.userRepo.GetUserByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetUserSaldo: failed to get user by ID")
		return 0, err
	}

	saldo, err = usecase.userRepo.GetUserSaldo(ctx, user.ID)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetUserSaldo: failed to get user saldo by ID")
		return 0, err
	}

	return saldo, nil
}

func (usecase *userUsecase) UpdateUserSaldo(ctx context.Context, id string, newSaldo int64) error {
	user, err := usecase.userRepo.GetUserByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("usecase.UpdateUserSaldo: failed to get user by ID")
		return err
	}
	err = usecase.userRepo.UpdateUserSaldo(ctx, user.ID, newSaldo)
	if err != nil {
		logrus.WithError(err).Error("usecase.UpdateUserSaldo: failed to update user saldo")
		return err
	}
	return nil
}
