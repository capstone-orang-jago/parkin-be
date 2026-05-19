package handler

import (
	"net/http"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type lokasiParkirHandler struct {
	lokasiParkirUsecase domain.LokasiParkirUsecase
}

func NewLokasiParkirHandler(lokasiParkirUsecase domain.LokasiParkirUsecase) domain.LokasiParkirHandler {
	return &lokasiParkirHandler{
		lokasiParkirUsecase: lokasiParkirUsecase,
	}
}

// CreateLokasiParkir godoc
// @Summary Create a new lokasi parkir
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.LokasiParkirRequestDTO true "Create Lokasi Parkir Request"
// @Success 201 {object} dto.LokasiParkirResponseDTO
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /admin/lokasi-parkir [post]
func (handler *lokasiParkirHandler) CreateLokasiParkir(c *gin.Context) {
	var req dto.LokasiParkirRequestDTO
	if err := helpers.DecodeJSON(c, &req); err != nil {
		logrus.WithError(err).Error("handler.CreateLokasiParkir: failed to decode request body")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}

	lokasiParkir, err := handler.lokasiParkirUsecase.CreateLokasiParkir(c.Request.Context(), &req)
	if err != nil {
		logrus.WithError(err).Error("handler.CreateLokasiParkir: failed to create lokasi parkir")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusCreated,
		Message: "lokasi parkir created successfully",
		Data:    lokasiParkir,
	})
}

// GetAllLokasiParkir godoc
// @Summary Get all lokasi parkir
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.LokasiParkirResponseDTO
// @Failure 500 {object} domain.HttpResponse
// @Router /admin/lokasi-parkir [get]
func (handler *lokasiParkirHandler) GetAllLokasiParkir(c *gin.Context) {
	lokasiParkir, err := handler.lokasiParkirUsecase.GetAllLokasiParkir(c.Request.Context())
	if err != nil {
		logrus.WithError(err).Error("handler.GetAllLokasiParkir: failed to get all lokasi parkir")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    lokasiParkir,
	})
}

// GetLokasiParkirByID godoc
// @Summary Get lokasi parkir by ID
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param id path string true "Lokasi Parkir ID"
// @Success 200 {object} dto.LokasiParkirResponseDTO
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /admin/lokasi-parkir/{id} [get]
func (handler *lokasiParkirHandler) GetLokasiParkirByID(c *gin.Context) {
	id := c.Param("id")
	lokasiParkir, err := handler.lokasiParkirUsecase.GetLokasiParkirByID(c.Request.Context(), id)
	if err != nil {
		logrus.WithError(err).Error("handler.GetLokasiParkirByID: failed to get lokasi parkir by ID")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    lokasiParkir,
	})
}

// GetLokasiParkirByPetugasID godoc
// @Summary Get lokasi parkir by Petugas ID
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param petugas_id path string true "Petugas ID"
// @Success 200 {object} dto.LokasiParkirResponseDTO
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /admin/lokasi-parkir/petugas/{petugas_id} [get]
func (handler *lokasiParkirHandler) GetLokasiParkirByPetugasID(c *gin.Context) {
	petugasID := c.Param("petugas_id")
	lokasiParkir, err := handler.lokasiParkirUsecase.GetLokasiParkirByPetugasID(c.Request.Context(), petugasID)
	if err != nil {
		logrus.WithError(err).Error("handler.GetLokasiParkirByPetugasID: failed to get lokasi parkir by petugas ID")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    lokasiParkir,
	})
}

// UpdateLokasiParkir godoc
// @Summary Update lokasi parkir
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Lokasi Parkir ID"
// @Param request body dto.LokasiParkirRequestDTO true "Update Lokasi Parkir Request"
// @Success 200 {object} dto.LokasiParkirResponseDTO
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /admin/lokasi-parkir/{id} [put]
func (handler *lokasiParkirHandler) UpdateLokasiParkir(c *gin.Context) {
	id := c.Param("id")

	var req dto.LokasiParkirRequestDTO
	if err := helpers.DecodeJSON(c, &req); err != nil {
		logrus.WithError(err).Error("handler.UpdateLokasiParkir: failed to decode request body")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}

	lokasiParkir, err := handler.lokasiParkirUsecase.UpdateLokasiParkir(c.Request.Context(), req, id)
	if err != nil {
		logrus.WithError(err).Error("handler.UpdateLokasiParkir: failed to update lokasi parkir")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "lokasi parkir updated successfully",
		Data:    lokasiParkir,
	})
}

// DeleteLokasiParkir godoc
// @Summary Delete lokasi parkir
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param id path string true "Lokasi Parkir ID"
// @Success 200 {object} domain.HttpResponse
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /admin/lokasi-parkir/{id} [delete]
func (handler *lokasiParkirHandler) DeleteLokasiParkir(c *gin.Context) {
	id := c.Param("id")

	err := handler.lokasiParkirUsecase.DeleteLokasiParkir(c.Request.Context(), id)
	if err != nil {
		logrus.WithError(err).Error("handler.DeleteLokasiParkir: failed to delete lokasi parkir")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "lokasi parkir deleted successfully",
	})
}
