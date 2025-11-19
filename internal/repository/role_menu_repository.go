package repository

import (
	"ipl-be-svc/internal/models"

	"gorm.io/gorm"
)

// RoleMenuRepository defines the interface for role menu data operations
type RoleMenuRepository interface {
	Create(roleMenu *models.RoleMenu) error
	GetByID(id uint) (*models.RoleMenu, error)
	GetAll(limit, offset int) ([]models.RoleMenu, int64, error)
	Update(roleMenu *models.RoleMenu) error
	Delete(id uint) error
	GetByRoleID(roleID uint) ([]models.RoleMenu, error)
	AttachMasterMenu(roleMenuID, masterMenuID uint, order *float64) error
	DetachMasterMenu(roleMenuID, masterMenuID uint) error
	AttachRole(roleMenuID, roleID uint, order *float64) error
	DetachRole(roleMenuID, roleID uint) error
	GetWithRelations(id uint) (*models.RoleMenu, error)
}

// roleMenuRepository implements RoleMenuRepository
type roleMenuRepository struct {
	db *gorm.DB
}

// NewRoleMenuRepository creates a new instance of RoleMenuRepository
func NewRoleMenuRepository(db *gorm.DB) RoleMenuRepository {
	return &roleMenuRepository{
		db: db,
	}
}

// Create creates a new role menu
func (r *roleMenuRepository) Create(roleMenu *models.RoleMenu) error {
	return r.db.Create(roleMenu).Error
}

// GetByID retrieves a role menu by ID
func (r *roleMenuRepository) GetByID(id uint) (*models.RoleMenu, error) {
	var roleMenu models.RoleMenu
	err := r.db.First(&roleMenu, id).Error
	if err != nil {
		return nil, err
	}
	return &roleMenu, nil
}

// GetWithRelations retrieves a role menu by ID with its master menus and roles
func (r *roleMenuRepository) GetWithRelations(id uint) (*models.RoleMenu, error) {
	var roleMenu models.RoleMenu
	err := r.db.Preload("MasterMenus").Preload("Roles").First(&roleMenu, id).Error
	if err != nil {
		return nil, err
	}
	return &roleMenu, nil
}

// GetAll retrieves all role menus with pagination
func (r *roleMenuRepository) GetAll(limit, offset int) ([]models.RoleMenu, int64, error) {
	var roleMenus []models.RoleMenu
	var total int64

	// Count total records
	if err := r.db.Model(&models.RoleMenu{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records with preloaded relations
	query := r.db.Preload("MasterMenus").Preload("Roles").Order("role_menu_ord ASC, id ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&roleMenus).Error; err != nil {
		return nil, 0, err
	}

	return roleMenus, total, nil
}

// Update updates a role menu
func (r *roleMenuRepository) Update(roleMenu *models.RoleMenu) error {
	return r.db.Save(roleMenu).Error
}

// Delete deletes a role menu by ID
func (r *roleMenuRepository) Delete(id uint) error {
	// Delete associations first
	r.db.Exec("DELETE FROM role_menus_master_menu_lnk WHERE role_menu_id = ?", id)
	r.db.Exec("DELETE FROM role_menus_role_lnk WHERE role_menu_id = ?", id)

	// Delete the role menu
	return r.db.Delete(&models.RoleMenu{}, id).Error
}

// GetByRoleID retrieves role menus associated with a specific role
func (r *roleMenuRepository) GetByRoleID(roleID uint) ([]models.RoleMenu, error) {
	var roleMenus []models.RoleMenu
	err := r.db.Joins("JOIN role_menus_role_lnk ON role_menus.id = role_menus_role_lnk.role_menu_id").
		Where("role_menus_role_lnk.role_id = ?", roleID).
		Preload("MasterMenus").
		Order("role_menus_role_lnk.role_menu_ord ASC, role_menus.id ASC").
		Find(&roleMenus).Error
	return roleMenus, err
}

// AttachMasterMenu attaches a master menu to a role menu
func (r *roleMenuRepository) AttachMasterMenu(roleMenuID, masterMenuID uint, order *float64) error {
	link := models.RoleMenuMasterMenuLink{
		RoleMenuID:   roleMenuID,
		MasterMenuID: masterMenuID,
		RoleMenuOrd:  order,
	}
	return r.db.Create(&link).Error
}

// DetachMasterMenu detaches a master menu from a role menu
func (r *roleMenuRepository) DetachMasterMenu(roleMenuID, masterMenuID uint) error {
	return r.db.Where("role_menu_id = ? AND master_menu_id = ?", roleMenuID, masterMenuID).
		Delete(&models.RoleMenuMasterMenuLink{}).Error
}

// AttachRole attaches a role to a role menu
func (r *roleMenuRepository) AttachRole(roleMenuID, roleID uint, order *float64) error {
	link := models.RoleMenuRoleLink{
		RoleMenuID:  roleMenuID,
		RoleID:      roleID,
		RoleMenuOrd: order,
	}
	return r.db.Create(&link).Error
}

// DetachRole detaches a role from a role menu
func (r *roleMenuRepository) DetachRole(roleMenuID, roleID uint) error {
	return r.db.Where("role_menu_id = ? AND role_id = ?", roleMenuID, roleID).
		Delete(&models.RoleMenuRoleLink{}).Error
}
