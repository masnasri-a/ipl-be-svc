package handler

import (
	"strconv"

	"ipl-be-svc/internal/service"
	"ipl-be-svc/pkg/logger"
	"ipl-be-svc/pkg/utils"

	"github.com/gin-gonic/gin"
)

// RoleMenuHandler handles role menu-related HTTP requests
type RoleMenuHandler struct {
	roleMenuService service.RoleMenuService
	logger          *logger.Logger
}

// NewRoleMenuHandler creates a new role menu handler
func NewRoleMenuHandler(roleMenuService service.RoleMenuService, logger *logger.Logger) *RoleMenuHandler {
	return &RoleMenuHandler{
		roleMenuService: roleMenuService,
		logger:          logger,
	}
}

// CreateRoleMenu handles POST /api/v1/role-menus
// @Summary Create a new role menu
// @Description Create a new role menu with the provided information
// @Tags role-menus
// @Accept json
// @Produce json
// @Param role_menu body service.CreateRoleMenuRequest true "Role menu data"
// @Success 201 {object} utils.APIResponse{data=models.RoleMenu} "Role menu created successfully"
// @Failure 400 {object} utils.APIResponse "Invalid request data"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/role-menus [post]
func (h *RoleMenuHandler) CreateRoleMenu(c *gin.Context) {
	var req service.CreateRoleMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid create role menu request")
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	roleMenu, err := h.roleMenuService.CreateRoleMenu(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create role menu")
		utils.InternalServerErrorResponse(c, "Failed to create role menu", err)
		return
	}

	h.logger.WithField("id", roleMenu.ID).Info("Role menu created successfully")
	utils.CreatedResponse(c, "Role menu created successfully", roleMenu)
}

// GetRoleMenu handles GET /api/v1/role-menus/:id
// @Summary Get role menu by ID
// @Description Get role menu information by ID with master menus and roles
// @Tags role-menus
// @Accept json
// @Produce json
// @Param id path int true "Role Menu ID"
// @Success 200 {object} utils.APIResponse{data=models.RoleMenu} "Role menu retrieved successfully"
// @Failure 400 {object} utils.APIResponse "Invalid role menu ID"
// @Failure 404 {object} utils.APIResponse "Role menu not found"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/role-menus/{id} [get]
func (h *RoleMenuHandler) GetRoleMenu(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.WithError(err).WithField("id_param", idParam).Error("Invalid role menu ID parameter")
		utils.BadRequestResponse(c, "Invalid role menu ID", err)
		return
	}

	roleMenu, err := h.roleMenuService.GetRoleMenuByID(uint(id))
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to get role menu")

		if err.Error() == "record not found" {
			utils.NotFoundResponse(c, "Role menu not found")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to get role menu", err)
		return
	}

	utils.SuccessResponse(c, "Role menu retrieved successfully", roleMenu)
}

// GetAllRoleMenus handles GET /api/v1/role-menus
// @Summary Get all role menus
// @Description Get all role menus with pagination and relations
// @Tags role-menus
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]models.RoleMenu} "Role menus retrieved successfully"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/role-menus [get]
func (h *RoleMenuHandler) GetAllRoleMenus(c *gin.Context) {
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

	roleMenus, total, err := h.roleMenuService.GetAllRoleMenus(limit, offset)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get role menus")
		utils.InternalServerErrorResponse(c, "Failed to get role menus", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "Role menus retrieved successfully", roleMenus, page, limit, total)
}

// UpdateRoleMenu handles PUT /api/v1/role-menus/:id
// @Summary Update role menu
// @Description Update role menu information
// @Tags role-menus
// @Accept json
// @Produce json
// @Param id path int true "Role Menu ID"
// @Param role_menu body service.UpdateRoleMenuRequest true "Role menu update data"
// @Success 200 {object} utils.APIResponse{data=models.RoleMenu} "Role menu updated successfully"
// @Failure 400 {object} utils.APIResponse "Invalid request data"
// @Failure 404 {object} utils.APIResponse "Role menu not found"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/role-menus/{id} [put]
func (h *RoleMenuHandler) UpdateRoleMenu(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.WithError(err).WithField("id_param", idParam).Error("Invalid role menu ID parameter")
		utils.BadRequestResponse(c, "Invalid role menu ID", err)
		return
	}

	var req service.UpdateRoleMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid update role menu request")
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	roleMenu, err := h.roleMenuService.UpdateRoleMenu(uint(id), &req)
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to update role menu")

		if err.Error() == "record not found" {
			utils.NotFoundResponse(c, "Role menu not found")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to update role menu", err)
		return
	}

	h.logger.WithField("id", id).Info("Role menu updated successfully")
	utils.SuccessResponse(c, "Role menu updated successfully", roleMenu)
}

// DeleteRoleMenu handles DELETE /api/v1/role-menus/:id
// @Summary Delete role menu
// @Description Delete role menu by ID and its associations
// @Tags role-menus
// @Accept json
// @Produce json
// @Param id path int true "Role Menu ID"
// @Success 200 {object} utils.APIResponse "Role menu deleted successfully"
// @Failure 400 {object} utils.APIResponse "Invalid role menu ID"
// @Failure 404 {object} utils.APIResponse "Role menu not found"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/role-menus/{id} [delete]
func (h *RoleMenuHandler) DeleteRoleMenu(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.WithError(err).WithField("id_param", idParam).Error("Invalid role menu ID parameter")
		utils.BadRequestResponse(c, "Invalid role menu ID", err)
		return
	}

	err = h.roleMenuService.DeleteRoleMenu(uint(id))
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to delete role menu")

		if err.Error() == "record not found" {
			utils.NotFoundResponse(c, "Role menu not found")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to delete role menu", err)
		return
	}

	h.logger.WithField("id", id).Info("Role menu deleted successfully")
	utils.SuccessResponse(c, "Role menu deleted successfully", nil)
}

// GetRoleMenusByRoleID handles GET /api/v1/roles/:role_id/role-menus
// @Summary Get role menus by role ID
// @Description Get all role menus associated with a specific role
// @Tags role-menus
// @Accept json
// @Produce json
// @Param role_id path int true "Role ID"
// @Success 200 {object} utils.APIResponse{data=[]models.RoleMenu} "Role menus retrieved successfully"
// @Failure 400 {object} utils.APIResponse "Invalid role ID"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/roles/{role_id}/role-menus [get]
func (h *RoleMenuHandler) GetRoleMenusByRoleID(c *gin.Context) {
	roleIDParam := c.Param("role_id")
	roleID, err := strconv.ParseUint(roleIDParam, 10, 32)
	if err != nil {
		h.logger.WithError(err).WithField("role_id_param", roleIDParam).Error("Invalid role ID parameter")
		utils.BadRequestResponse(c, "Invalid role ID", err)
		return
	}

	roleMenus, err := h.roleMenuService.GetRoleMenusByRoleID(uint(roleID))
	if err != nil {
		h.logger.WithError(err).WithField("role_id", roleID).Error("Failed to get role menus by role ID")
		utils.InternalServerErrorResponse(c, "Failed to get role menus", err)
		return
	}

	utils.SuccessResponse(c, "Role menus retrieved successfully", roleMenus)
}

// AttachMasterMenu handles POST /api/v1/role-menus/:id/master-menus
// @Summary Attach master menu to role menu
// @Description Attach a master menu to a role menu with optional ordering
// @Tags role-menus
// @Accept json
// @Produce json
// @Param id path int true "Role Menu ID"
// @Param attach_request body service.AttachMasterMenuRequest true "Attach master menu data"
// @Success 200 {object} utils.APIResponse "Master menu attached successfully"
// @Failure 400 {object} utils.APIResponse "Invalid request data"
// @Failure 404 {object} utils.APIResponse "Role menu or master menu not found"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/role-menus/{id}/master-menus [post]
func (h *RoleMenuHandler) AttachMasterMenu(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.WithError(err).WithField("id_param", idParam).Error("Invalid role menu ID parameter")
		utils.BadRequestResponse(c, "Invalid role menu ID", err)
		return
	}

	var req service.AttachMasterMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid attach master menu request")
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	err = h.roleMenuService.AttachMasterMenuToRoleMenu(uint(id), req.MasterMenuID, req.Order)
	if err != nil {
		h.logger.WithError(err).WithFields(map[string]interface{}{
			"role_menu_id":   id,
			"master_menu_id": req.MasterMenuID,
		}).Error("Failed to attach master menu to role menu")

		if err.Error() == "role menu not found" || err.Error() == "master menu not found" {
			utils.NotFoundResponse(c, err.Error())
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to attach master menu", err)
		return
	}

	utils.SuccessResponse(c, "Master menu attached successfully", nil)
}

// DetachMasterMenu handles DELETE /api/v1/role-menus/:id/master-menus/:master_menu_id
// @Summary Detach master menu from role menu
// @Description Detach a master menu from a role menu
// @Tags role-menus
// @Accept json
// @Produce json
// @Param id path int true "Role Menu ID"
// @Param master_menu_id path int true "Master Menu ID"
// @Success 200 {object} utils.APIResponse "Master menu detached successfully"
// @Failure 400 {object} utils.APIResponse "Invalid ID"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/role-menus/{id}/master-menus/{master_menu_id} [delete]
func (h *RoleMenuHandler) DetachMasterMenu(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.WithError(err).WithField("id_param", idParam).Error("Invalid role menu ID parameter")
		utils.BadRequestResponse(c, "Invalid role menu ID", err)
		return
	}

	masterMenuIDParam := c.Param("master_menu_id")
	masterMenuID, err := strconv.ParseUint(masterMenuIDParam, 10, 32)
	if err != nil {
		h.logger.WithError(err).WithField("master_menu_id_param", masterMenuIDParam).Error("Invalid master menu ID parameter")
		utils.BadRequestResponse(c, "Invalid master menu ID", err)
		return
	}

	err = h.roleMenuService.DetachMasterMenuFromRoleMenu(uint(id), uint(masterMenuID))
	if err != nil {
		h.logger.WithError(err).WithFields(map[string]interface{}{
			"role_menu_id":   id,
			"master_menu_id": masterMenuID,
		}).Error("Failed to detach master menu from role menu")
		utils.InternalServerErrorResponse(c, "Failed to detach master menu", err)
		return
	}

	utils.SuccessResponse(c, "Master menu detached successfully", nil)
}

// AttachRole handles POST /api/v1/role-menus/:id/roles
// @Summary Attach role to role menu
// @Description Attach a role to a role menu with optional ordering
// @Tags role-menus
// @Accept json
// @Produce json
// @Param id path int true "Role Menu ID"
// @Param attach_request body service.AttachRoleRequest true "Attach role data"
// @Success 200 {object} utils.APIResponse "Role attached successfully"
// @Failure 400 {object} utils.APIResponse "Invalid request data"
// @Failure 404 {object} utils.APIResponse "Role menu not found"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/role-menus/{id}/roles [post]
func (h *RoleMenuHandler) AttachRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.WithError(err).WithField("id_param", idParam).Error("Invalid role menu ID parameter")
		utils.BadRequestResponse(c, "Invalid role menu ID", err)
		return
	}

	var req service.AttachRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid attach role request")
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	err = h.roleMenuService.AttachRoleToRoleMenu(uint(id), req.RoleID, req.Order)
	if err != nil {
		h.logger.WithError(err).WithFields(map[string]interface{}{
			"role_menu_id": id,
			"role_id":      req.RoleID,
		}).Error("Failed to attach role to role menu")

		if err.Error() == "role menu not found" {
			utils.NotFoundResponse(c, err.Error())
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to attach role", err)
		return
	}

	utils.SuccessResponse(c, "Role attached successfully", nil)
}

// DetachRole handles DELETE /api/v1/role-menus/:id/roles/:role_id
// @Summary Detach role from role menu
// @Description Detach a role from a role menu
// @Tags role-menus
// @Accept json
// @Produce json
// @Param id path int true "Role Menu ID"
// @Param role_id path int true "Role ID"
// @Success 200 {object} utils.APIResponse "Role detached successfully"
// @Failure 400 {object} utils.APIResponse "Invalid ID"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/role-menus/{id}/roles/{role_id} [delete]
func (h *RoleMenuHandler) DetachRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.WithError(err).WithField("id_param", idParam).Error("Invalid role menu ID parameter")
		utils.BadRequestResponse(c, "Invalid role menu ID", err)
		return
	}

	roleIDParam := c.Param("role_id")
	roleID, err := strconv.ParseUint(roleIDParam, 10, 32)
	if err != nil {
		h.logger.WithError(err).WithField("role_id_param", roleIDParam).Error("Invalid role ID parameter")
		utils.BadRequestResponse(c, "Invalid role ID", err)
		return
	}

	err = h.roleMenuService.DetachRoleFromRoleMenu(uint(id), uint(roleID))
	if err != nil {
		h.logger.WithError(err).WithFields(map[string]interface{}{
			"role_menu_id": id,
			"role_id":      roleID,
		}).Error("Failed to detach role from role menu")
		utils.InternalServerErrorResponse(c, "Failed to detach role", err)
		return
	}

	utils.SuccessResponse(c, "Role detached successfully", nil)
}
