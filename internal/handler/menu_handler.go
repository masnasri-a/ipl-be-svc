package handler

import (
	"strings"

	"github.com/gin-gonic/gin"

	"ipl-be-svc/internal/service"
	"ipl-be-svc/pkg/utils"
)

// MenuHandler handles menu-related HTTP requests
type MenuHandler struct {
	menuService service.MenuService
}

// NewMenuHandler creates a new menu handler
func NewMenuHandler(menuService service.MenuService) *MenuHandler {
	return &MenuHandler{
		menuService: menuService,
	}
}

// MenuResponse represents the menu response structure
type MenuResponse struct {
	ID          uint    `json:"id" example:"1"`
	DocumentID  string  `json:"document_id" example:"mo5qqs8ezbruui07t91p6da8"`
	NamaMenu    string  `json:"nama_menu" example:"Master Data"`
	KodeMenu    string  `json:"kode_menu" example:"master-data"`
	UrutanMenu  *int    `json:"urutan_menu" example:"1"`
	IsActive    *bool   `json:"is_active" example:"true"`
	PublishedAt *string `json:"published_at,omitempty" example:"2025-10-23T15:16:28.206Z"`
}

// GetMenusByUserID handles GET /api/v1/menus/user/:user_id
// @Summary Get menus by user ID
// @Description Get list of menus accessible by a specific user ID
// @Tags menus
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} utils.APIResponse{data=[]MenuResponse} "Menus retrieved successfully"
// @Failure 400 {object} utils.APIResponse "Invalid user ID"
// @Failure 404 {object} utils.APIResponse "No menus found"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/menus/user/{user_id} [get]
func (h *MenuHandler) GetMenusByUserID(c *gin.Context) {
	userID, err := utils.GetIDParam(c)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID", err)
		return
	}

	menus, err := h.menuService.GetMenusByUserID(userID)
	if err != nil {
		if strings.Contains(err.Error(), "invalid user ID") {
			utils.BadRequestResponse(c, "Invalid user ID", err)
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to get menus", err)
		return
	}

	if len(menus) == 0 {
		utils.SuccessResponse(c, "No menus found for this user", []interface{}{})
		return
	}

	// Convert to response format
	var menuResponses []MenuResponse
	for _, menu := range menus {
		var publishedAt *string
		if menu.PublishedAt != nil {
			pubAt := menu.PublishedAt.Format("2006-01-02T15:04:05.000Z")
			publishedAt = &pubAt
		}

		menuResponses = append(menuResponses, MenuResponse{
			ID:          menu.ID,
			DocumentID:  menu.DocumentID,
			NamaMenu:    menu.NamaMenu,
			KodeMenu:    menu.KodeMenu,
			UrutanMenu:  menu.UrutanMenu,
			IsActive:    menu.IsActive,
			PublishedAt: publishedAt,
		})
	}

	utils.SuccessResponse(c, "Menus retrieved successfully", menuResponses)
}