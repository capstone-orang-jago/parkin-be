package usecase

import (
	"context"
	"errors"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/sirupsen/logrus"
)

type lokasiParkirUsecase struct {
	lokasiParkirRepo domain.LokasiParkirRepository
}

func NewLokasiParkirUsecase(lokasiParkirRepo domain.LokasiParkirRepository) domain.LokasiParkirUsecase {
	return &lokasiParkirUsecase{
		lokasiParkirRepo: lokasiParkirRepo}
}

func (usecase *lokasiParkirUsecase) CreateLokasiParkir(ctx context.Context, req *dto.LokasiParkirRequestDTO) (res *dto.LokasiParkirResponseDTO, err error) {
	if req.NamaLokasi == "" {
		return nil, errors.New("nama lokasi are required")
	}
	if req.Alamat == "" {
		return nil, errors.New("alamat are required")
	}
	if req.PetugasID == "" {
		return nil, errors.New("petugas id are required")
	}
	if !req.IsActive {
		return nil, errors.New("is active are required")
	}

	lokasiParkir := &domain.LokasiParkir{
		NamaLokasi: req.NamaLokasi,
		Alamat:     req.Alamat,
		PetugasID:  req.PetugasID,
		IsActive:   req.IsActive,
	}

	err = usecase.lokasiParkirRepo.CreateLokasiParkir(ctx, lokasiParkir)
	if err != nil {
		logrus.WithError(err).Error("usecase.CreateLokasiParkir: failed to create lokasi parkir")
		return nil, err
	}

	response := &dto.LokasiParkirResponseDTO{
		ID:         lokasiParkir.ID,
		NamaLokasi: lokasiParkir.NamaLokasi,
		Alamat:     lokasiParkir.Alamat,
		PetugasID:  lokasiParkir.PetugasID,
		IsActive:   lokasiParkir.IsActive,
		CreatedAt:  lokasiParkir.CreatedAt,
		UpdatedAt:  lokasiParkir.UpdatedAt,
	}

	return response, err
}

func (usecase *lokasiParkirUsecase) GetAllLokasiParkir(ctx context.Context) ([]*dto.LokasiParkirResponseDTO, error) {
	lokasiParkir, err := usecase.lokasiParkirRepo.GetAllLokasiParkir(ctx)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetAllLokasiParkir: failed to get all lokasi parkir")
		return nil, err
	}

	responses := make([]*dto.LokasiParkirResponseDTO, 0, len(lokasiParkir))
	for _, lokasi := range lokasiParkir {
		response := &dto.LokasiParkirResponseDTO{
			ID:         lokasi.ID,
			NamaLokasi: lokasi.NamaLokasi,
			Alamat:     lokasi.Alamat,
			PetugasID:  lokasi.PetugasID,
			IsActive:   lokasi.IsActive,
			CreatedAt:  lokasi.CreatedAt,
			UpdatedAt:  lokasi.UpdatedAt,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (usecase *lokasiParkirUsecase) GetLokasiParkirByID(ctx context.Context, id string) (*dto.LokasiParkirResponseDTO, error) {
	lokasiParkir, err := usecase.lokasiParkirRepo.GetLokasiParkirByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetLokasiParkirByID: failed to get lokasi parkir by id")
		return nil, err
	}

	response := &dto.LokasiParkirResponseDTO{
		ID:         lokasiParkir.ID,
		NamaLokasi: lokasiParkir.NamaLokasi,
		Alamat:     lokasiParkir.Alamat,
		PetugasID:  lokasiParkir.PetugasID,
		IsActive:   lokasiParkir.IsActive,
		CreatedAt:  lokasiParkir.CreatedAt,
		UpdatedAt:  lokasiParkir.UpdatedAt,
	}

	return response, err
}

func (usecase *lokasiParkirUsecase) GetLokasiParkirByPetugasID(ctx context.Context, petugasID string) (*dto.LokasiParkirResponseDTO, error) {
	lokasiParkir, err := usecase.lokasiParkirRepo.GetLokasiParkirByPetugasID(ctx, petugasID)
	if err != nil {
		logrus.WithError(err).Error("usecase.GetLokasiParkirByPetugasID: failed to get lokasi parkir by petugas id")
		return nil, err
	}

	response := &dto.LokasiParkirResponseDTO{
		ID:         lokasiParkir.ID,
		NamaLokasi: lokasiParkir.NamaLokasi,
		Alamat:     lokasiParkir.Alamat,
		PetugasID:  lokasiParkir.PetugasID,
		IsActive:   lokasiParkir.IsActive,
		CreatedAt:  lokasiParkir.CreatedAt,
		UpdatedAt:  lokasiParkir.UpdatedAt,
	}

	return response, err
}

func (usecase *lokasiParkirUsecase) UpdateLokasiParkir(ctx context.Context, req dto.LokasiParkirRequestDTO, id string) (res *dto.LokasiParkirResponseDTO, err error) {
	lokasiParkir, err := usecase.lokasiParkirRepo.GetLokasiParkirByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("usecase.UpdateLokasiParkir: failed to get lokasi parkir by id")
		return nil, err
	}

	if req.NamaLokasi == "" && req.Alamat == "" && req.PetugasID == "" && !req.IsActive {
		return nil, errors.New("at least one field must be provided to update")
	}

	if req.NamaLokasi == "" {
		lokasiParkir.NamaLokasi = req.NamaLokasi
	}
	if req.Alamat == "" {
		lokasiParkir.Alamat = req.Alamat
	}
	if req.PetugasID == "" {
		lokasiParkir.PetugasID = req.PetugasID
	}
	if !req.IsActive {
		lokasiParkir.IsActive = req.IsActive
	}

	err = usecase.lokasiParkirRepo.UpdateLokasiParkir(ctx, lokasiParkir, id)
	if err != nil {
		logrus.WithError(err).Error("usecase.UpdateLokasiParkir: failed to update lokasi parkir")
		return nil, err
	}

	response := &dto.LokasiParkirResponseDTO{
		ID:         lokasiParkir.ID,
		NamaLokasi: lokasiParkir.NamaLokasi,
		Alamat:     lokasiParkir.Alamat,
		PetugasID:  lokasiParkir.PetugasID,
		IsActive:   lokasiParkir.IsActive,
		CreatedAt:  lokasiParkir.CreatedAt,
		UpdatedAt:  lokasiParkir.UpdatedAt,
	}

	return response, err
}

func (usecase *lokasiParkirUsecase) DeleteLokasiParkir(ctx context.Context, id string) error {
	_, err := usecase.lokasiParkirRepo.GetLokasiParkirByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("usecase.DeleteLokasiParkir: failed to get lokasi parkir by id")
		return err
	}

	err = usecase.lokasiParkirRepo.DeleteLokasiParkir(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("usecase.DeleteLokasiParkir: failed to delete lokasi parkir")
		return err
	}

	return nil
}
