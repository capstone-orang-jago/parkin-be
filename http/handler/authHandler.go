package handler

import (
	"net/http"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type authHandler struct {
	authUsecase domain.AuthUsecase
	userUsecase domain.UserUsecase
}

func NewAuthHandler(authUsecase domain.AuthUsecase, userUsecase domain.UserUsecase) domain.AuthHandler {
	return &authHandler{
		authUsecase: authUsecase,
		userUsecase: userUsecase,
	}
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequestDTO true "Register Request"
// @Success 201 {object} dto.RegisterResponseDTO
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /auth/register [post]
func (handler *authHandler) Register(c *gin.Context) {
	var req dto.RegisterRequestDTO
	if err := helpers.DecodeJSON(c, &req); err != nil {
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}

	user, err := handler.authUsecase.Register(c.Request.Context(), req)
	if err != nil {
		logrus.WithError(err).Error("handler.Register: failed to register user")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	response := dto.RegisterResponseDTO{
		UserID:    user.UserID,
		Telephone: user.Telephone,
		OTPSentAt: user.OTPSentAt,
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusCreated,
		Message: "user created successfully",
		Data:    response,
	})
}

// VerifyRegister godoc
// @Summary Verify user registration
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyOtpCodeRequestDTO true "Verify OTP Code Request"
// @Success 200 {object} domain.Auth
// @Failure 400 {object} domain.HttpResponse
// @Failure 401 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /auth/verify-otp [post]
func (handler *authHandler) VerifyRegister(c *gin.Context) {
	var req dto.VerifyOtpCodeRequestDTO
	if err := helpers.DecodeJSON(c, &req); err != nil {
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}

	res, err := handler.authUsecase.VerifyRegister(c.Request.Context(), req)
	if err != nil {
		logrus.WithError(err).Error("handler.VerifyRegister: failed to verify registration")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "registration verified successfully",
		Data:    res,
	})
}

// Login godoc
// @Summary Login a user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequestDTO true "Login Request"
// @Success 200 {object} dto.LoginResponseDTO
// @Failure 400 {object} domain.HttpResponse
// @Failure 401 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /auth/login [post]
func (handler *authHandler) Login(c *gin.Context) {
	var req dto.LoginRequestDTO
	if err := helpers.DecodeJSON(c, &req); err != nil {
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}

	res, err := handler.authUsecase.Login(c.Request.Context(), req)
	if err != nil {
		logrus.WithError(err).Error("handler.Login: failed to login")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	}

	user, err := handler.userUsecase.GetUserByID(c.Request.Context(), res.UserID)
	if err != nil {
		logrus.WithError(err).Error("handler.Login: failed to get user by ID")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to get user by ID",
		})
		return
	}

	response := dto.LoginResponseDTO{
		UserID:       user.ID,
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "login successful",
		Data:    response,
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequestDTO true "Refresh Token Request"
// @Success 200 {object} domain.Auth
// @Failure 400 {object} domain.HttpResponse
// @Failure 401 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /auth/refresh [post]
func (handler *authHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequestDTO
	if err := helpers.DecodeJSON(c, &req); err != nil {
		logrus.WithError(err).Error("handler.RefreshToken: failed to decode request body")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}

	res, err := handler.authUsecase.RefreshToken(
		c.Request.Context(),
		req,
	)

	if err != nil {
		logrus.WithError(err).Error("handler.RefreshToken: failed to refresh token")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "refresh token successful",
		Data:    res,
	})
}

// Logout godoc
// @Summary Logout a user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequestDTO true "Refresh Token Request"
// @Security BearerAuth
// @Success 200 {object} domain.HttpResponse
// @Failure 400 {object} domain.HttpResponse
// @Failure 401 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /auth/logout [post]
func (handler *authHandler) Logout(c *gin.Context) {
	var req dto.RefreshTokenRequestDTO
	if err := helpers.DecodeJSON(c, &req); err != nil {
		logrus.WithError(err).Error("handler.Logout: failed to decode request body")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}

	userId := c.GetString("sub")

	err := handler.authUsecase.Logout(c.Request.Context(), userId, req.RefreshToken)
	if err != nil {
		logrus.WithError(err).Error("handler.Logout: failed to logout")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "logout successful",
	})
}
