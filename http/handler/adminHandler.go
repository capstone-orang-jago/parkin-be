package handler

import (
	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type adminHandler struct {
	adminUsecase domain.AdminUsecase
}

func NewAdminHandler(adminUsecase domain.AdminUsecase) domain.AdminHandler {
	return &adminHandler{adminUsecase: adminUsecase}
}

// GetAllTransactions godoc
// @Summary      Global Transaction Monitoring
// @Description  Retrieve global parking logs with rich analytical multi-filters for admin monitoring
// @Tags         Admin Transactions
// @Accept       json
// @Produce      json
// @Param        date_from    query     string  false  "Format: YYYY-MM-DD"
// @Param        date_to      query     string  false  "Format: YYYY-MM-DD"
// @Param        location_id  query     string  false  "Filter by Spot UUID"
// @Param        officer_id   query     string  false  "Filter by Officer UUID"
// @Param        vehicle      query     string  false  "motor / mobil"
// @Param        status       query     string  false  "sukses / gagal"
// @Param        page         query     int     false  "Page indicator"
// @Param        limit        query     int     false  "Row window scale"
// @Security     BearerAuth
// @Success      200          {object}  domain.HttpResponse{data=dto.PaginatedTransactionResponse}
// @Router       /admin/transactions [get]
func (h *adminHandler) GetAllTransactions(c *gin.Context) {
	var query domain.AdminTransaksiQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		helpers.Response(c, domain.HttpResponse{Code: 400, Message: err.Error()})
		return
	}

	data, total, err := h.adminUsecase.GetAllTransactions(c.Request.Context(), query)
	if err != nil {
		logrus.WithError(err).Error("handler.Admin.GetAllTransactions: failed")
		helpers.Response(c, domain.HttpResponse{Code: 500, Message: "failed to load transactions"})
		return
	}

	helpers.Response(c, domain.HttpResponse{
		Code:    200,
		Message: "global logs retrieved successfully",
		Data: gin.H{
			"total": total,
			"page":  query.Page,
			"limit": query.Limit,
			"items": data,
		},
	})
}

// GetTransactionDetail godoc
// @Summary      Transaction Audit Audit (Masked ID)
// @Description  Get precise logs data with strict masked user security patterns
// @Tags         Admin Transactions
// @Param        id   path      string  true  "Transaction UUID"
// @Security     BearerAuth
// @Router       /admin/transactions/{id} [get]
func (h *adminHandler) GetTransactionDetail(c *gin.Context) {
	id := c.Param("id")
	data, err := h.adminUsecase.GetTransactionDetail(c.Request.Context(), id)
	if err != nil {
		helpers.Response(c, domain.HttpResponse{Code: 500, Message: "failed to load transaction detail"})
		return
	}
	helpers.Response(c, domain.HttpResponse{Code: 200, Message: "success", Data: data})
}

// GetFilterOptions godoc
// @Summary      Get Transaction Filter Options
// @Description  Retrieve dynamic dropdown options for admin transaction monitoring filter (locations, officers, vehicles, statuses)
// @Tags         Admin Transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  domain.HttpResponse{data=dto.AdminFilterOptionsResponseDTO}
// @Router       /admin/transactions/filter-options [get]
func (h *adminHandler) GetFilterOptions(c *gin.Context) {
	data, err := h.adminUsecase.GetFilterOptions(c.Request.Context())
	if err != nil {
		logrus.WithError(err).Error("handler.Admin.GetFilterOptions: failed")
		helpers.Response(c, domain.HttpResponse{Code: 500, Message: "failed to load filter options"})
		return
	}
	helpers.Response(c, domain.HttpResponse{Code: 200, Message: "filter options retrieved", Data: data})
}

// GetAllLocations godoc
// @Summary      Get All Parking Locations
// @Description  Retrieve a list of all parking spots/locations with their assigned officers
// @Tags         Admin Locations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  domain.HttpResponse{data=[]dto.LokasiParkirResponseDTO}
// @Router       /admin/locations [get]
func (h *adminHandler) GetAllLocations(c *gin.Context) {
	data, err := h.adminUsecase.GetAllLocations(c.Request.Context())
	if err != nil {
		logrus.WithError(err).Error("handler.Admin.GetAllLocations: failed")
		helpers.Response(c, domain.HttpResponse{Code: 500, Message: "failed to load locations"})
		return
	}
	helpers.Response(c, domain.HttpResponse{Code: 200, Message: "locations retrieved successfully", Data: data})
}

// GetLocationDetail godoc
// @Summary      Get Location Detail & Stats
// @Description  Retrieve detailed statistics and recent transaction history for a specific parking spot
// @Tags         Admin Locations
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Location UUID"
// @Security     BearerAuth
// @Success      200  {object}  domain.HttpResponse{data=dto.AdminLocationDetailDTO}
// @Router       /admin/locations/{id} [get]
func (h *adminHandler) GetLocationDetail(c *gin.Context) {
	id := c.Param("id")
	data, err := h.adminUsecase.GetLocationDetail(c.Request.Context(), id)
	if err != nil {
		logrus.WithError(err).Error("handler.Admin.GetLocationDetail: failed")
		helpers.Response(c, domain.HttpResponse{Code: 404, Message: err.Error()})
		return
	}
	helpers.Response(c, domain.HttpResponse{Code: 200, Message: "location detail retrieved", Data: data})
}

// GetAllOfficers godoc
// @Summary      Get All Field Officers
// @Description  Retrieve a list of all parking officers/dishub personnel
// @Tags         Admin Officers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  domain.HttpResponse{data=[]dto.UserResponseDTO}
// @Router       /admin/officers [get]
func (h *adminHandler) GetAllOfficers(c *gin.Context) {
	data, err := h.adminUsecase.GetAllOfficers(c.Request.Context())
	if err != nil {
		logrus.WithError(err).Error("handler.Admin.GetAllOfficers: failed")
		helpers.Response(c, domain.HttpResponse{Code: 500, Message: "failed to load officers"})
		return
	}
	helpers.Response(c, domain.HttpResponse{Code: 200, Message: "officers retrieved successfully", Data: data})
}

// GetOfficerDetail godoc
// @Summary      Get Officer Performance Detail
// @Description  Retrieve specific officer's duty performance and recent handled transactions
// @Tags         Admin Officers
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Officer UUID"
// @Security     BearerAuth
// @Success      200  {object}  domain.HttpResponse{data=dto.AdminOfficerDetailDTO}
// @Router       /admin/officers/{id} [get]
func (h *adminHandler) GetOfficerDetail(c *gin.Context) {
	id := c.Param("id")
	data, err := h.adminUsecase.GetOfficerDetail(c.Request.Context(), id)
	if err != nil {
		logrus.WithError(err).Error("handler.Admin.GetOfficerDetail: failed")
		helpers.Response(c, domain.HttpResponse{Code: 404, Message: err.Error()})
		return
	}
	helpers.Response(c, domain.HttpResponse{Code: 200, Message: "officer detail retrieved", Data: data})
}

// GetRevenueReport godoc
// @Summary      Get Aggregate Revenue Report
// @Description  Retrieve high-level financial reports with optional filters for business insights
// @Tags         Admin Reports
// @Accept       json
// @Produce      json
// @Param        date_from    query     string  false  "Format: YYYY-MM-DD"
// @Param        date_to      query     string  false  "Format: YYYY-MM-DD"
// @Param        location_id  query     string  false  "Filter by specific spot"
// @Param        vehicle      query     string  false  "motor / mobil"
// @Security     BearerAuth
// @Success      200  {object}  domain.HttpResponse{data=dto.AdminReportResponseDTO}
// @Router       /admin/reports [get]
func (h *adminHandler) GetRevenueReport(c *gin.Context) {
	var query domain.AdminReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		logrus.WithError(err).Error("handler.Admin.GetRevenueReport: invalid query params")
		helpers.Response(c, domain.HttpResponse{Code: 400, Message: "invalid query parameters"})
		return
	}

	data, err := h.adminUsecase.GetRevenueReport(c.Request.Context(), query)
	if err != nil {
		logrus.WithError(err).Error("handler.Admin.GetRevenueReport: failed")
		helpers.Response(c, domain.HttpResponse{Code: 500, Message: "failed to generate revenue report"})
		return
	}
	helpers.Response(c, domain.HttpResponse{Code: 200, Message: "revenue report generated", Data: data})
}
