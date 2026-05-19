package handler

import (
	"net/http"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/gin-gonic/gin"
)

type midtransHandler struct {
	midtransUsecase domain.MidtransUsecase
}

func NewMidtransHandler(midtransUsecase domain.MidtransUsecase) domain.MidtransHandler {
	return &midtransHandler{
		midtransUsecase: midtransUsecase,
	}
}

// CreateTopUp godoc
// @Summary Initiate a new top-up transaction
// @Description Initiate a new top-up transaction with Midtrans and get the payment URL
// @Tags user
// @Accept json
// @Produce json
// @Param request body dto.TopUpRequestDTO true "Top-Up Request"
// @Security BearerAuth
// @Success 200 {object} dto.TopUpResponseDTO
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /topup/create [post]
func (handler *midtransHandler) CreateTopUp(c *gin.Context) {
	userID := c.GetString("sub")

	var req dto.TopUpRequestDTO
	if err := helpers.DecodeJSON(c, &req); err != nil {
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request body: " + err.Error(),
		})
		return
	}

	res, err := handler.midtransUsecase.InitiateTopUp(c.Request.Context(), userID, req)
	if err != nil {
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to initiate top-up: " + err.Error(),
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "top-up initiated successfully",
		Data:    res,
	})
}

// MidtransCallback godoc
// @Summary Handle Midtrans payment callback
// @Description Handle the payment callback from Midtrans and update top-up status accordingly
// @Tags midtrans
// @Accept json
// @Produce json
// @Param request body dto.MidtransCallbackDTO true "Midtrans Callback Request"
// @Success 200 {object} domain.HttpResponse
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /midtrans/callback [post]
func (handler *midtransHandler) MidtransCallback(c *gin.Context) {
	var callbackReq dto.MidtransCallbackDTO
	if err := c.ShouldBindBodyWithJSON(&callbackReq); err != nil {
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid callback request body: " + err.Error(),
		})
		return
	}

	success, err := handler.midtransUsecase.TopUpSaldo(c.Request.Context(), callbackReq)
	if err != nil || !success {
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to update top-up status",
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "top-up status updated successfully",
	})
}

// GetTopUpStatus godoc
// @Summary Get the status of a top-up transaction
// @Description Get the status of a top-up transaction by its order ID
// @Tags user
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Security BearerAuth
// @Success 200 {object} domain.HttpResponse
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /topup/{order_id}/status [get]
func (handler *midtransHandler) GetTopUpStatus(c *gin.Context) {
	orderID := c.Param("order_id")
	if orderID == "" {
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "order id is required",
		})
		return
	}

	topUpData, err := handler.midtransUsecase.GetTopUpStatus(c.Request.Context(), orderID)
	if err != nil {
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to fetch top-up status: " + err.Error(),
		})
		return
	}

	if topUpData == nil {
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusNotFound,
			Message: "top-up record not found",
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "top-up status fetched successfully",
		Data:    topUpData,
	})
}
