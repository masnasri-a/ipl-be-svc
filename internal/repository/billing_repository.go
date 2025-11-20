package repository

import (
	"ipl-be-svc/internal/models"

	"gorm.io/gorm"
)

// BillingRepository defines the interface for billing data operations
type BillingRepository interface {
	GetBillingByID(id uint) (*models.Billing, error)
	GetUsersWithPenghuniRole() ([]*models.User, error)
	GetActiveMonthlySettingBillings() ([]*models.SettingBilling, error)
	CreateBulkBillings(billings []*models.Billing) error
	CreateBulkBillingProfileLinks(links []*models.BillingProfileLink) error
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

// GetUsersWithPenghuniRole retrieves all users with role type "penghuni"
func (r *billingRepository) GetUsersWithPenghuniRole() ([]*models.User, error) {
	var users []*models.User

	err := r.db.Table("up_users").
		Joins("JOIN up_users_role_lnk url ON up_users.id = url.user_id").
		Joins("JOIN up_roles r ON url.role_id = r.id").
		Where("r.type = ?", "penghuni").
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetActiveMonthlySettingBillings retrieves all active monthly setting billings
func (r *billingRepository) GetActiveMonthlySettingBillings() ([]*models.SettingBilling, error) {
	var settings []*models.SettingBilling

	err := r.db.Where("jenis_billing = ? AND is_active = ? AND published_at IS NOT NULL", "bulanan", true).Find(&settings).Error
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// CreateBulkBillings creates multiple billing records in a transaction
func (r *billingRepository) CreateBulkBillings(billings []*models.Billing) error {
	return r.db.CreateInBatches(billings, 100).Error
}

// CreateBulkBillingProfileLinks creates multiple billing-profile links in a transaction
func (r *billingRepository) CreateBulkBillingProfileLinks(links []*models.BillingProfileLink) error {
	return r.db.CreateInBatches(links, 100).Error
}
