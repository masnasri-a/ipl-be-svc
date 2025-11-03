package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel contains common columns for all tables
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// MasterMenu represents the master_menus table
type MasterMenu struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	DocumentID  string    `json:"document_id" gorm:"column:document_id"`
	NamaMenu    string    `json:"nama_menu" gorm:"column:nama_menu"`
	KodeMenu    string    `json:"kode_menu" gorm:"column:kode_menu"`
	UrutanMenu  *int      `json:"urutan_menu" gorm:"column:urutan_menu"`
	IsActive    *bool     `json:"is_active" gorm:"column:is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
	Locale      *string   `json:"locale"`
}

// TableName sets the insert table name for MasterMenu
func (MasterMenu) TableName() string {
	return "master_menus"
}

// Billing represents the billings table
type Billing struct {
	ID          uint       `json:"id" gorm:"primarykey"`
	DocumentID  *string    `json:"document_id" gorm:"column:document_id"`
	Bulan       *int       `json:"bulan" gorm:"column:bulan"`
	Tahun       *int       `json:"tahun" gorm:"column:tahun"`
	Nominal     *int64     `json:"nominal" gorm:"column:nominal"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedByID *int       `json:"created_by_id"`
	UpdatedByID *int       `json:"updated_by_id"`
	Locale      *string    `json:"locale"`
}

// TableName sets the insert table name for Billing
func (Billing) TableName() string {
	return "billings"
}