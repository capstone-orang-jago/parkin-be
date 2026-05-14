package handler

import (
	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/gin-gonic/gin"
)

type otpCodeHandler struct {
	otpCodeUsecase domain.OtpCodeUsecase
}

func NewOtpCodeHandler(otpCodeUsecase domain.OtpCodeUsecase) domain.OtpCodeHandler {
	return &otpCodeHandler{
		otpCodeUsecase: otpCodeUsecase,
	}
}

// ResendOtpCode godoc
// @Summary Resend OTP code to user's telephone
// @Tags otp
// @Accept json
// @Produce json
// @Param request body dto.OtpCodeRequestDTO true "OTP Code Request"
// @Success 200 {object} domain.HttpResponse
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /otp/resend [post]
func (h *otpCodeHandler) ResendOtpCode(c *gin.Context) {
	var req dto.OtpCodeRequestDTO
	if err := helpers.DecodeJSON(c, &req); err != nil {
		helpers.Response(c, domain.HttpResponse{
			Code:    400,
			Message: "invalid request body",
		})
		return
	}

	if err := h.otpCodeUsecase.SendOtpCode(c.Request.Context(), req.Telephone); err != nil {
		helpers.Response(c, domain.HttpResponse{
			Code:    500,
			Message: "failed to send OTP code",
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    200,
		Message: "OTP code sent successfully",
	})
}
