package handler

import (
	"net/http"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type userHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(userUsecase domain.UserUsecase) domain.UserHandler {
	return &userHandler{
		userUsecase: userUsecase,
	}
}

// GetAllUsers godoc
// @Summary Get all users
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.UserResponseDTO
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /user [get]
func (handler *userHandler) GetAllUsers(c *gin.Context) {
	users, err := handler.userUsecase.GetAllUsers(c.Request.Context())
	if err != nil {
		logrus.WithError(err).Error("handler.GetAllUsers: failed to get all users")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    users,
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponseDTO
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /user/{id} [get]
func (handler *userHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := handler.userUsecase.GetUserByID(c.Request.Context(), id)
	if err != nil {
		logrus.WithError(err).Error("handler.GetUserByID: failed to get user by ID")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    user,
	})
}

// GetUserProfile godoc
// @Summary Get user profile
// @Tags profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.UserResponseDTO
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /user/profile [get]
// @Router /officer/profile [get]
func (handler *userHandler) GetUserProfile(c *gin.Context) {
	userId := c.GetString("sub")
	role := c.GetString("role")

	profile, err := handler.userUsecase.GetUserProfile(c.Request.Context(), userId, role)
	if err != nil {
		logrus.WithError(err).Error("handler.GetUserProfile: failed to get user profile")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    profile,
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.UserRequestDTO true "User Request"
// @Security BearerAuth
// @Success 201 {object} dto.UserResponseDTO
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /user [post]
func (handler *userHandler) CreateUser(c *gin.Context) {
	var req dto.UserRequestDTO
	if err := helpers.DecodeJSON(c, &req); err != nil {
		logrus.WithError(err).Error("handler.CreateUser: failed to decode request body")
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

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusCreated,
		Message: "user created successfully",
		Data:    user,
	})
}

// UpdateUser godoc
// @Summary Update user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body dto.UserRequestDTO true "User Request"
// @Security BearerAuth
// @Success 200 {object} dto.UserResponseDTO
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /user/{id} [put]
func (handler *userHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req dto.UserRequestDTO
	if err := helpers.DecodeJSON(c, &req); err != nil {
		logrus.WithError(err).Error("handler.UpdateUser: failed to decode request body")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}

	user, err := handler.userUsecase.UpdateUser(c.Request.Context(), req, id)
	if err != nil {
		logrus.WithError(err).Error("handler.UpdateUser: failed to update user")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "user updated successfully",
		Data:    user,
	})
}

// DeleteUser godoc
// @Summary Delete user
// @Tags users
// @Param id path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} domain.HttpResponse
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /user/{id} [delete]
func (handler *userHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := handler.userUsecase.DeleteUser(c.Request.Context(), id)
	if err != nil {
		logrus.WithError(err).Error("handler.DeleteUser: failed to delete user")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "user deleted successfully",
	})
}

// GetUserSaldo godoc
// @Summary Get user saldo
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.HttpResponse{data=int64}
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /user/me/saldo [get]
func (handler *userHandler) GetUserSaldo(c *gin.Context) {
	userId := c.GetString("sub")

	saldo, err := handler.userUsecase.GetUserSaldo(c.Request.Context(), userId)
	if err != nil {
		logrus.WithError(err).Error("handler.GetUserSaldo: failed to get user saldo")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	response := dto.UserSaldoResponseDTO{
		ID:    userId,
		Saldo: saldo,
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

// DeductUserSaldo godoc
// @Summary Deduct user saldo
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.UserSaldoRequestDTO true "User Saldo Request"
// @Security BearerAuth
// @Success 200 {object} domain.HttpResponse
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Router /user/me/saldo [put]
func (handler *userHandler) DeductUserSaldo(c *gin.Context) {
	userId := c.GetString("sub")

	var req dto.UserSaldoRequestDTO
	if err := helpers.DecodeJSON(c, &req); err != nil {
		logrus.WithError(err).Error("handler.DeductUserSaldo: failed to decode request body")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}

	err := handler.userUsecase.DeductUserSaldo(c.Request.Context(), userId, req.Nominal)
	if err != nil {
		logrus.WithError(err).Error("handler.DeductUserSaldo: failed to deduct user saldo")
		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "user saldo deducted successfully",
	})
}
