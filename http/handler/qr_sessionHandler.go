package handler

import (
	"net/http"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/gin-gonic/gin"
)

type qrSessionHandler struct {
	qrSessionUsecase domain.QRSessionUsecase
}

func NewQRSessionHandler(qrSessionUsecase domain.QRSessionUsecase) domain.QRSessionHandler {
	return &qrSessionHandler{
		qrSessionUsecase: qrSessionUsecase,
	}
}

// CreateQRSession godoc
// @Summary Create a new QR session for petugas
// @Tags qr-session
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 201 {object} dto.QRSessionResponseDTO
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /officer/qr-code [post]
func (handler *qrSessionHandler) CreateQRSession(c *gin.Context) {
	userID := c.GetString("sub")

	qrSession, err := handler.qrSessionUsecase.CreateQRSession(c.Request.Context(), userID)
	if err != nil {
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to create QR session",
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusCreated,
		Message: "QR session created successfully",
		Data:    qrSession,
	})
}
