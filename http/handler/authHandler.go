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
// @Param request body dto.UserRequestDTO true "User Request"
// @Success 201 {object} dto.UserResponseDTO
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /auth/register [post]
func (handler *authHandler) Register(c *gin.Context) {
	var req dto.UserRequestDTO
	if err := helpers.DecodeJSON(c, &req); err != nil {
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}

	user, err := handler.userUsecase.CreateUser(c.Request.Context(), req)
	if err != nil {
		logrus.WithError(err).Error("handler.CreateUser: failed to create user")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	response := dto.UserResponseDTO{
		ID:         user.ID,
		Name:       user.Name,
		Telephone:  user.Telephone,
		Email:      user.Email,
		Role:       user.Role,
		Saldo:      user.Saldo,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusCreated,
		Message: "user created successfully",
		Data:    response,
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
// @Router /auth/refresh-token [post]
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

	err := handler.authUsecase.Logout(c.Request.Context(), req.RefreshToken)
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
