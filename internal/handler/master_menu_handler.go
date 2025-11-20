package handler

import (
	"strconv"

	"ipl-be-svc/internal/service"
	"ipl-be-svc/pkg/logger"
	"ipl-be-svc/pkg/utils"

	"github.com/gin-gonic/gin"
)

// MasterMenuHandler handles master menu-related HTTP requests
type MasterMenuHandler struct {
	masterMenuService service.MasterMenuService
	logger            *logger.Logger
}

// NewMasterMenuHandler creates a new master menu handler
func NewMasterMenuHandler(masterMenuService service.MasterMenuService, logger *logger.Logger) *MasterMenuHandler {
	return &MasterMenuHandler{
		masterMenuService: masterMenuService,
		logger:            logger,
	}
}

// CreateMasterMenu handles POST /api/v1/master-menus
// @Summary Create a new master menu
// @Description Create a new master menu with the provided information
// @Tags master-menus
// @Accept json
// @Produce json
// @Param master_menu body service.CreateMasterMenuRequest true "Master menu data"
// @Success 201 {object} utils.APIResponse{data=models.MasterMenu} "Master menu created successfully"
// @Failure 400 {object} utils.APIResponse "Invalid request data"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/master-menus [post]
func (h *MasterMenuHandler) CreateMasterMenu(c *gin.Context) {
	var req service.CreateMasterMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid create master menu request")
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	masterMenu, err := h.masterMenuService.CreateMasterMenu(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create master menu")
		utils.InternalServerErrorResponse(c, "Failed to create master menu", err)
		return
	}

	h.logger.WithField("id", masterMenu.ID).Info("Master menu created successfully")
	utils.CreatedResponse(c, "Master menu created successfully", masterMenu)
}

// GetMasterMenu handles GET /api/v1/master-menus/:id
// @Summary Get master menu by ID
// @Description Get master menu information by ID
// @Tags master-menus
// @Accept json
// @Produce json
// @Param id path int true "Master Menu ID"
// @Success 200 {object} utils.APIResponse{data=models.MasterMenu} "Master menu retrieved successfully"
// @Failure 400 {object} utils.APIResponse "Invalid master menu ID"
// @Failure 404 {object} utils.APIResponse "Master menu not found"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/master-menus/{id} [get]
func (h *MasterMenuHandler) GetMasterMenu(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.WithError(err).WithField("id_param", idParam).Error("Invalid master menu ID parameter")
		utils.BadRequestResponse(c, "Invalid master menu ID", err)
		return
	}

	masterMenu, err := h.masterMenuService.GetMasterMenuByID(uint(id))
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to get master menu")

		if err.Error() == "record not found" {
			utils.NotFoundResponse(c, "Master menu not found")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to get master menu", err)
		return
	}

	utils.SuccessResponse(c, "Master menu retrieved successfully", masterMenu)
}

// GetAllMasterMenus handles GET /api/v1/master-menus
// @Summary Get all master menus
// @Description Get all master menus with pagination
// @Tags master-menus
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]models.MasterMenu} "Master menus retrieved successfully"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/master-menus [get]
func (h *MasterMenuHandler) GetAllMasterMenus(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	masterMenus, total, err := h.masterMenuService.GetAllMasterMenus(limit, offset)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get master menus")
		utils.InternalServerErrorResponse(c, "Failed to get master menus", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "Master menus retrieved successfully", masterMenus, page, limit, total)
}

// UpdateMasterMenu handles PUT /api/v1/master-menus/:id
// @Summary Update master menu
// @Description Update master menu information
// @Tags master-menus
// @Accept json
// @Produce json
// @Param id path int true "Master Menu ID"
// @Param master_menu body service.UpdateMasterMenuRequest true "Master menu update data"
// @Success 200 {object} utils.APIResponse{data=models.MasterMenu} "Master menu updated successfully"
// @Failure 400 {object} utils.APIResponse "Invalid request data"
// @Failure 404 {object} utils.APIResponse "Master menu not found"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/master-menus/{id} [put]
func (h *MasterMenuHandler) UpdateMasterMenu(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.WithError(err).WithField("id_param", idParam).Error("Invalid master menu ID parameter")
		utils.BadRequestResponse(c, "Invalid master menu ID", err)
		return
	}

	var req service.UpdateMasterMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid update master menu request")
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	masterMenu, err := h.masterMenuService.UpdateMasterMenu(uint(id), &req)
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to update master menu")

		if err.Error() == "record not found" {
			utils.NotFoundResponse(c, "Master menu not found")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to update master menu", err)
		return
	}

	h.logger.WithField("id", id).Info("Master menu updated successfully")
	utils.SuccessResponse(c, "Master menu updated successfully", masterMenu)
}

// DeleteMasterMenu handles DELETE /api/v1/master-menus/:id
// @Summary Delete master menu
// @Description Delete master menu by ID
// @Tags master-menus
// @Accept json
// @Produce json
// @Param id path int true "Master Menu ID"
// @Success 200 {object} utils.APIResponse "Master menu deleted successfully"
// @Failure 400 {object} utils.APIResponse "Invalid master menu ID"
// @Failure 404 {object} utils.APIResponse "Master menu not found"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/master-menus/{id} [delete]
func (h *MasterMenuHandler) DeleteMasterMenu(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.WithError(err).WithField("id_param", idParam).Error("Invalid master menu ID parameter")
		utils.BadRequestResponse(c, "Invalid master menu ID", err)
		return
	}

	err = h.masterMenuService.DeleteMasterMenu(uint(id))
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to delete master menu")

		if err.Error() == "record not found" {
			utils.NotFoundResponse(c, "Master menu not found")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to delete master menu", err)
		return
	}

	h.logger.WithField("id", id).Info("Master menu deleted successfully")
	utils.SuccessResponse(c, "Master menu deleted successfully", nil)
}
