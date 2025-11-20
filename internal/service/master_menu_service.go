package service

import (
	"fmt"
	"ipl-be-svc/internal/models"
	"ipl-be-svc/internal/repository"
	"ipl-be-svc/pkg/logger"
)

// MasterMenuService interface defines master menu service methods
type MasterMenuService interface {
	CreateMasterMenu(req *CreateMasterMenuRequest) (*models.MasterMenu, error)
	GetMasterMenuByID(id uint) (*models.MasterMenu, error)
	GetAllMasterMenus(limit, offset int) ([]models.MasterMenu, int64, error)
	UpdateMasterMenu(id uint, req *UpdateMasterMenuRequest) (*models.MasterMenu, error)
	DeleteMasterMenu(id uint) error
}

// CreateMasterMenuRequest represents the request to create a master menu
type CreateMasterMenuRequest struct {
	DocumentID *string `json:"document_id" example:"menu001"`
	NamaMenu   string  `json:"nama_menu" binding:"required" example:"Dashboard"`
	KodeMenu   string  `json:"kode_menu" binding:"required" example:"DASHBOARD"`
	UrutanMenu *int    `json:"urutan_menu" example:"1"`
	IsActive   *bool   `json:"is_active" example:"true"`
	Locale     *string `json:"locale" example:"id"`
}

// UpdateMasterMenuRequest represents the request to update a master menu
type UpdateMasterMenuRequest struct {
	DocumentID *string `json:"document_id" example:"menu001"`
	NamaMenu   *string `json:"nama_menu" example:"Dashboard"`
	KodeMenu   *string `json:"kode_menu" example:"DASHBOARD"`
	UrutanMenu *int    `json:"urutan_menu" example:"1"`
	IsActive   *bool   `json:"is_active" example:"true"`
	Locale     *string `json:"locale" example:"id"`
}

// masterMenuService implements MasterMenuService interface
type masterMenuService struct {
	masterMenuRepo repository.MasterMenuRepository
	logger         *logger.Logger
}

// NewMasterMenuService creates a new master menu service
func NewMasterMenuService(masterMenuRepo repository.MasterMenuRepository, logger *logger.Logger) MasterMenuService {
	return &masterMenuService{
		masterMenuRepo: masterMenuRepo,
		logger:         logger,
	}
}

// CreateMasterMenu creates a new master menu
func (s *masterMenuService) CreateMasterMenu(req *CreateMasterMenuRequest) (*models.MasterMenu, error) {
	// Validate required fields
	if req.NamaMenu == "" {
		return nil, fmt.Errorf("nama_menu is required")
	}
	if req.KodeMenu == "" {
		return nil, fmt.Errorf("kode_menu is required")
	}

	// Check if kode_menu already exists
	existing, _ := s.masterMenuRepo.GetByKodeMenu(req.KodeMenu)
	if existing != nil {
		return nil, fmt.Errorf("kode_menu already exists")
	}

	// Create master menu
	masterMenu := &models.MasterMenu{
		NamaMenu:   req.NamaMenu,
		KodeMenu:   req.KodeMenu,
		UrutanMenu: req.UrutanMenu,
		IsActive:   req.IsActive,
		Locale:     req.Locale,
	}

	if req.DocumentID != nil {
		masterMenu.DocumentID = *req.DocumentID
	}

	err := s.masterMenuRepo.Create(masterMenu)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create master menu")
		return nil, err
	}

	s.logger.WithFields(map[string]interface{}{
		"id":        masterMenu.ID,
		"kode_menu": masterMenu.KodeMenu,
		"nama_menu": masterMenu.NamaMenu,
	}).Info("Master menu created successfully")

	return masterMenu, nil
}

// GetMasterMenuByID retrieves a master menu by ID
func (s *masterMenuService) GetMasterMenuByID(id uint) (*models.MasterMenu, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid master menu ID")
	}

	masterMenu, err := s.masterMenuRepo.GetByID(id)
	if err != nil {
		s.logger.WithError(err).WithField("id", id).Error("Failed to get master menu")
		return nil, err
	}

	return masterMenu, nil
}

// GetAllMasterMenus retrieves all master menus with pagination
func (s *masterMenuService) GetAllMasterMenus(limit, offset int) ([]models.MasterMenu, int64, error) {
	masterMenus, total, err := s.masterMenuRepo.GetAll(limit, offset)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get master menus")
		return nil, 0, err
	}

	return masterMenus, total, nil
}

// UpdateMasterMenu updates a master menu
func (s *masterMenuService) UpdateMasterMenu(id uint, req *UpdateMasterMenuRequest) (*models.MasterMenu, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid master menu ID")
	}

	// Get existing master menu
	masterMenu, err := s.masterMenuRepo.GetByID(id)
	if err != nil {
		s.logger.WithError(err).WithField("id", id).Error("Failed to get master menu for update")
		return nil, err
	}

	// Update fields if provided
	if req.DocumentID != nil {
		masterMenu.DocumentID = *req.DocumentID
	}
	if req.NamaMenu != nil {
		masterMenu.NamaMenu = *req.NamaMenu
	}
	if req.KodeMenu != nil {
		// Check if kode_menu already exists (excluding current record)
		existing, _ := s.masterMenuRepo.GetByKodeMenu(*req.KodeMenu)
		if existing != nil && existing.ID != id {
			return nil, fmt.Errorf("kode_menu already exists")
		}
		masterMenu.KodeMenu = *req.KodeMenu
	}
	if req.UrutanMenu != nil {
		masterMenu.UrutanMenu = req.UrutanMenu
	}
	if req.IsActive != nil {
		masterMenu.IsActive = req.IsActive
	}
	if req.Locale != nil {
		masterMenu.Locale = req.Locale
	}

	err = s.masterMenuRepo.Update(masterMenu)
	if err != nil {
		s.logger.WithError(err).WithField("id", id).Error("Failed to update master menu")
		return nil, err
	}

	s.logger.WithFields(map[string]interface{}{
		"id":        masterMenu.ID,
		"kode_menu": masterMenu.KodeMenu,
		"nama_menu": masterMenu.NamaMenu,
	}).Info("Master menu updated successfully")

	return masterMenu, nil
}

// DeleteMasterMenu deletes a master menu
func (s *masterMenuService) DeleteMasterMenu(id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid master menu ID")
	}

	// Check if master menu exists
	_, err := s.masterMenuRepo.GetByID(id)
	if err != nil {
		s.logger.WithError(err).WithField("id", id).Error("Master menu not found for deletion")
		return err
	}

	err = s.masterMenuRepo.Delete(id)
	if err != nil {
		s.logger.WithError(err).WithField("id", id).Error("Failed to delete master menu")
		return err
	}

	s.logger.WithField("id", id).Info("Master menu deleted successfully")
	return nil
}
