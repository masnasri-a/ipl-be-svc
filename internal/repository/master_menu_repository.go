package repository

import (
	"ipl-be-svc/internal/models"

	"gorm.io/gorm"
)

// MasterMenuRepository defines the interface for master menu data operations
type MasterMenuRepository interface {
	Create(masterMenu *models.MasterMenu) error
	GetByID(id uint) (*models.MasterMenu, error)
	GetAll(limit, offset int) ([]models.MasterMenu, int64, error)
	Update(masterMenu *models.MasterMenu) error
	Delete(id uint) error
	GetByKodeMenu(kodeMenu string) (*models.MasterMenu, error)
}

// masterMenuRepository implements MasterMenuRepository
type masterMenuRepository struct {
	db *gorm.DB
}

// NewMasterMenuRepository creates a new instance of MasterMenuRepository
func NewMasterMenuRepository(db *gorm.DB) MasterMenuRepository {
	return &masterMenuRepository{
		db: db,
	}
}

// Create creates a new master menu
func (r *masterMenuRepository) Create(masterMenu *models.MasterMenu) error {
	return r.db.Create(masterMenu).Error
}

// GetByID retrieves a master menu by ID
func (r *masterMenuRepository) GetByID(id uint) (*models.MasterMenu, error) {
	var masterMenu models.MasterMenu
	err := r.db.First(&masterMenu, id).Error
	if err != nil {
		return nil, err
	}
	return &masterMenu, nil
}

// GetAll retrieves all master menus with pagination
func (r *masterMenuRepository) GetAll(limit, offset int) ([]models.MasterMenu, int64, error) {
	var masterMenus []models.MasterMenu
	var total int64

	// Count total records
	if err := r.db.Model(&models.MasterMenu{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	query := r.db.Order("urutan_menu ASC, id ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&masterMenus).Error; err != nil {
		return nil, 0, err
	}

	return masterMenus, total, nil
}

// Update updates a master menu
func (r *masterMenuRepository) Update(masterMenu *models.MasterMenu) error {
	return r.db.Save(masterMenu).Error
}

// Delete deletes a master menu by ID
func (r *masterMenuRepository) Delete(id uint) error {
	return r.db.Delete(&models.MasterMenu{}, id).Error
}

// GetByKodeMenu retrieves a master menu by kode_menu
func (r *masterMenuRepository) GetByKodeMenu(kodeMenu string) (*models.MasterMenu, error) {
	var masterMenu models.MasterMenu
	err := r.db.Where("kode_menu = ?", kodeMenu).First(&masterMenu).Error
	if err != nil {
		return nil, err
	}
	return &masterMenu, nil
}
