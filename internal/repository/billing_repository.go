package repository

import (
	"ipl-be-svc/internal/models"

	"gorm.io/gorm"
)

// BillingRepository defines the interface for billing data operations
type BillingRepository interface {
	GetBillingByID(id uint) (*models.Billing, error)
}

// billingRepository implements BillingRepository
type billingRepository struct {
	db *gorm.DB
}

// NewBillingRepository creates a new instance of BillingRepository
func NewBillingRepository(db *gorm.DB) BillingRepository {
	return &billingRepository{
		db: db,
	}
}

// GetBillingByID retrieves a billing record by ID
func (r *billingRepository) GetBillingByID(id uint) (*models.Billing, error) {
	var billing models.Billing
	
	err := r.db.Where("id = ?", id).First(&billing).Error
	if err != nil {
		return nil, err
	}

	return &billing, nil
}