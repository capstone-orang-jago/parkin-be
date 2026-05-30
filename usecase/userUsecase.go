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

func (usecase *userUsecase) GetAllUsers(ctx context.Context) (res []*dto.UserResponseDTO, err error) {
	users, err := usecase.userRepo.GetAllUsers(ctx)

	if err != nil {
		logrus.WithError(err).Error("usecase.GetAllUsers: failed to get all users")
		return nil, err
	} else if users == nil {
		logrus.WithError(err).Info("usecase.GetAllUsers: no users found")
		return nil, nil
	}

	responses := make([]*dto.UserResponseDTO, 0, len(users))

	for _, user := range users {
		response := &dto.UserResponseDTO{
			ID:         user.ID,
			Name:       user.Name,
			Telephone:  user.Telephone,
			Email:      user.Email,
			Role:       user.Role,
			Saldo:      user.Saldo,
			IsVerified: user.IsVerified,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (usecase *userUsecase) GetUserByID(ctx context.Context, id string) (res *dto.UserResponseDTO, err error) {
	user, err := usecase.userRepo.GetUserByID(ctx, id)

	if err != nil {
		logrus.WithError(err).Error("usecase.GetUserByID: failed to get user by ID")
		return nil, err
	} else if user == nil {
		logrus.WithError(err).Info("usecase.GetUserByID: no user found for ID")
		return nil, nil
	}

	response := &dto.UserResponseDTO{
		ID:         user.ID,
		Name:       user.Name,
		Telephone:  user.Telephone,
		Email:      user.Email,
		Role:       user.Role,
		Saldo:      user.Saldo,
		IsVerified: user.IsVerified,
	}

	return response, nil
}

func (usecase *userUsecase) GetUserProfile(ctx context.Context, id, role string) (res *dto.UserResponseDTO, err error) {
	user, err := usecase.userRepo.GetUserProfile(ctx, id, role)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetUserProfile: failed to get user profile")
		return nil, err
	} else if user == nil {
		logrus.WithError(err).Info("usecase.GetUserProfile: no user found for ID")
		return nil, nil
	}

	response := &dto.UserResponseDTO{
		ID:         user.ID,
		Name:       user.Name,
		Telephone:  user.Telephone,
		Email:      user.Email,
		Role:       user.Role,
		Saldo:      user.Saldo,
		IsVerified: user.IsVerified,
	}

	if role == "petugas" && user.LokasiParkir != nil {
		response.Lokasi = &dto.LokasiParkirResponseDTO{
			ID:         user.LokasiParkir.ID,
			NamaLokasi: user.LokasiParkir.NamaLokasi,
			Alamat:     user.LokasiParkir.Alamat,
			IsActive:   user.LokasiParkir.IsActive,
		}
	}

	return response, nil
}

func (usecase *userUsecase) CreateUser(ctx context.Context, req dto.UserRequestDTO) (res *dto.UserResponseDTO, err error) {
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
		Role:       "petugas",
		Saldo:      0,
		IsVerified: false,
	}

	err = usecase.userRepo.CreateUser(ctx, user)
	if err != nil {
		logrus.WithError(err).Error("usecase.CreateUser: failed to create user")
		return nil, err
	}

	response := &dto.UserResponseDTO{
		ID:         user.ID,
		Name:       user.Name,
		Telephone:  user.Telephone,
		Email:      user.Email,
		Role:       user.Role,
		Saldo:      user.Saldo,
		IsVerified: user.IsVerified,
	}

	return response, nil
}

func (usecase *userUsecase) UpdateUser(ctx context.Context, req dto.UserRequestDTO, id string) (res *dto.UserResponseDTO, err error) {
	user, err := usecase.userRepo.GetUserByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("usecase.UpdateUser: failed to get user by ID")
		return nil, err
	} else if user == nil {
		logrus.WithError(err).Info("usecase.UpdateUser: no user found for ID")
		return nil, nil
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

	response := &dto.UserResponseDTO{
		ID:         user.ID,
		Name:       user.Name,
		Telephone:  user.Telephone,
		Email:      user.Email,
		Role:       user.Role,
		Saldo:      user.Saldo,
		IsVerified: user.IsVerified,
	}

	return response, nil
}

func (usecase *userUsecase) DeleteUser(ctx context.Context, id string) error {
	user, err := usecase.userRepo.GetUserByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("usecase.DeleteUser: failed to get user by ID")
		return err
	} else if user == nil {
		logrus.WithError(err).Info("usecase.DeleteUser: no user found for ID")
		return nil
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
	} else if user == nil {
		logrus.WithError(err).Info("usecase.GetUserSaldo: no user found for ID")
		return 0, nil
	}

	saldo, err = usecase.userRepo.GetUserSaldo(ctx, user.ID)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetUserSaldo: failed to get user saldo by ID")
		return 0, err
	}

	return saldo, nil
}

func (usecase *userUsecase) DeductUserSaldo(ctx context.Context, id string, amount int64) error {
	user, err := usecase.userRepo.GetUserByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("usecase.DeductUserSaldo: failed to get user by ID")
		return err
	} else if user == nil {
		logrus.WithError(err).Info("usecase.DeductUserSaldo: no user found for ID")
		return nil
	}

	err = usecase.userRepo.DeductUserSaldo(ctx, user.ID, amount)
	if err != nil {
		logrus.WithError(err).Error("usecase.DeductUserSaldo: failed to deduct user saldo")
		return err
	}

	return nil
}
