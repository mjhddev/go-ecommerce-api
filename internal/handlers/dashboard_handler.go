package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mjhddev/go-ecommerce-api/internal/response"
	"github.com/mjhddev/go-ecommerce-api/internal/services"
)

type DashboardHandler struct {
	dashboardService services.DashboardService
}

func NewDashboardHandler(
	dashboardService services.DashboardService,
) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetDashboard godoc
//
//	@Summary		Get dashboard analytics
//	@Description	Get dashboard analytics for admin
//	@Tags			Admin Dashboard
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.SuccessResponse{data=dto.DashboardResponse}
//	@Failure		401	{object}	response.ErrorResponse
//	@Failure		403	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/admin/dashboard [get]
func (h *DashboardHandler) GetDashboard(c *gin.Context) {

	dashboard, err := h.dashboardService.GetDashboard()
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Dashboard retrieved successfully",
		dashboard,
	)
}
