package handler

import (
	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type transaksiParkirHandler struct {
	transaksiParkirUsecase domain.TransaksiParkirUsecase
}

func NewTransaksiParkirHandler(transaksiParkirUsecase domain.TransaksiParkirUsecase) domain.TransaksiParkirHandler {
	return &transaksiParkirHandler{transaksiParkirUsecase: transaksiParkirUsecase}
}

// CreateTransaksiParkir godoc
// @Summary Create a new parking transaction
// @Description Create a new parking transaction based on the provided QR payload and vehicle type
// @Tags Transaksi Parkir
// @Accept json
// @Produce json
// @Param request body dto.TransaksiParkirRequestDTO true "Transaksi Parkir Request"
// @Security BearerAuth
// @Success 201 {object} dto.TransaksiParkirResponseDTO
// @Router /transaksi-parkir [post]
func (handler *transaksiParkirHandler) CreateTransaksiParkir(c *gin.Context) {
	userID := c.GetString("sub")

	var req dto.TransaksiParkirRequestDTO

	if err := helpers.DecodeJSON(c, &req); err != nil {
		logrus.WithError(err).Error("handler.CreateTransaksiParkir: failed to decode request body")
		helpers.Response(c, domain.HttpResponse{
			Code:    400,
			Message: "invalid request body",
		})
		return
	}

	transaksiParkir, err := handler.transaksiParkirUsecase.CreateTransaksiParkir(c.Request.Context(), &req, userID)
	if err != nil {
		logrus.WithError(err).Error("handler.CreateTransaksiParkir: failed to create transaksi parkir")
		helpers.Response(c, domain.HttpResponse{
			Code:    500,
			Message: "failed to create transaksi parkir",
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    201,
		Message: "transaksi parkir created successfully",
		Data:    transaksiParkir,
	})
}

// GetTransaksiParkirByID godoc
// @Summary Get parking transaction by ID
// @Description Retrieve a parking transaction by its ID
// @Tags Transaksi Parkir
// @Accept json
// @Produce json
// @Param id path string true "Transaksi Parkir ID"
// @Security BearerAuth
// @Success 200 {object} dto.TransaksiParkirResponseDTO
// @Router /transaksi-parkir/{id} [get]
func (handler *transaksiParkirHandler) GetTransaksiParkirByID(c *gin.Context) {
	id := c.Param("id")

	transaksiParkir, err := handler.transaksiParkirUsecase.GetTransaksiParkirByID(c.Request.Context(), id)
	if err != nil {
		logrus.WithError(err).Error("handler.GetTransaksiParkirByID: failed to get transaksi parkir by ID")
		helpers.Response(c, domain.HttpResponse{
			Code:    500,
			Message: "failed to get transaksi parkir",
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    200,
		Message: "transaksi parkir retrieved successfully",
		Data:    transaksiParkir,
	})
}

// GetAllTransaksiParkirByUserID godoc
// @Summary Get all parking transactions for a user
// @Description Retrieve all parking transactions for the authenticated user, with optional pagination and filtering
// @Tags Transaksi Parkir
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page for pagination"
// @Param filter query string false "Filter for transactions (semua, minggu, bulan)"
// @Security BearerAuth
// @Success 200 {array} dto.TransaksiParkirResponseDTO
// @Router /user/transactions [get]
// @Router /officer/transactions [get]
func (handler *transaksiParkirHandler) GetAllTransaksiParkirByUserID(c *gin.Context) {
	userID := c.GetString("sub")
	role := c.GetString("role")

	var query domain.TransaksiQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		logrus.WithError(err).Error("handler.GetAllTransaksiParkirByUserID: failed to decode query parameters")
		helpers.Response(c, domain.HttpResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}
	if query.Filter == "" {
		query.Filter = "semua"
	}

	transaksiParkirList, err := handler.transaksiParkirUsecase.GetAllTransaksiParkirByUserID(c.Request.Context(), userID, role, query.Page, query.Limit, query.Filter)
	if err != nil {
		logrus.WithError(err).Error("handler.GetAllTransaksiParkirByUserID: failed to get all transaksi parkir by user ID")
		helpers.Response(c, domain.HttpResponse{
			Code:    500,
			Message: "failed to get transaksi parkir",
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    200,
		Message: "transaksi parkir retrieved successfully",
		Data:    transaksiParkirList,
	})
}
