package handler

import (
	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type dashboardHandler struct {
	dashboardUsecase domain.DashboardUsecase
}

func NewDashboardHandler(dashboardUsecase domain.DashboardUsecase) domain.DashboardHandler {
	return &dashboardHandler{dashboardUsecase: dashboardUsecase}
}

// GetOfficerDashboard godoc
// @Summary Get officer dashboard data
// @Description Retrieve dashboard data for the authenticated parking officer, including profile info and recent transactions
// @Tags officer
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.OfficerDashboardResponseDTO
// @Router /officer/dashboard [get]
func (handler *dashboardHandler) GetOfficerDashboard(c *gin.Context) {
	officerID := c.GetString("sub")

	dashboardData, err := handler.dashboardUsecase.GetOfficerDashboard(c.Request.Context(), officerID)
	if err != nil {
		logrus.WithError(err).Error("handler.GetOfficerDashboard: failed to get officer dashboard data")
		helpers.Response(c, domain.HttpResponse{
			Code:    500,
			Message: "failed to get dashboard data",
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    200,
		Message: "dashboard data retrieved successfully",
		Data:    dashboardData,
	})
}

// GetAdminDashboard godoc
// @Summary      Get Admin Dashboard Analytics
// @Description  Retrieve global analytics, trend metrics, recent streams, and spot performance using fast Go-runtime mapping
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  domain.HttpResponse{data=dto.AdminDashboardResponseDTO}
// @Router       /admin/dashboard [get]
func (handler *dashboardHandler) GetAdminDashboard(c *gin.Context) {
	dashboardData, err := handler.dashboardUsecase.GetAdminDashboard(c.Request.Context())
	if err != nil {
		logrus.WithError(err).Error("handler.GetAdminDashboard: failed to crunch admin dashboard analytics")
		helpers.Response(c, domain.HttpResponse{
			Code:    500,
			Message: "failed to crunch dashboard analytics",
		})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    200,
		Message: "admin dashboard analytics generated successfully",
		Data:    dashboardData,
	})
}
