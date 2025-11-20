package models

import (
	"time"
)

// MasterGeneralStatus represents the master_general_statuses table
type MasterGeneralStatus struct {
	ID                uint       `json:"id" gorm:"primarykey"`
	DocumentID        *string    `json:"document_id" gorm:"column:document_id"`
	Status            *string    `json:"status_name" gorm:"column:status_name"`
	StatusDescription *string    `json:"status_description" gorm:"column:status_description"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
	PublishedAt       *time.Time `json:"published_at"`
	CreatedByID       *int       `json:"created_by_id"`
	UpdatedByID       *int       `json:"updated_by_id"`
	Locale            *string    `json:"locale"`
}

// TableName sets the insert table name for MasterGeneralStatus
func (MasterGeneralStatus) TableName() string {
	return "master_general_statuses"
}
