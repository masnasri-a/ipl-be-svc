package models

import (
	"time"
)

// MasterKategoriTransaksi represents the master_kategori_transaksis table
type MasterKategoriTransaksi struct {
	ID          uint       `json:"id" gorm:"primarykey"`
	DocumentID  *string    `json:"document_id" gorm:"column:document_id"`
	Nama        *string    `json:"nama" gorm:"column:nama"`
	Keterangan  *string    `json:"keterangan" gorm:"column:keterangan"`
	Order       *int       `json:"order" gorm:"column:order"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedByID *int       `json:"created_by_id"`
	UpdatedByID *int       `json:"updated_by_id"`
	Locale      *string    `json:"locale"`
}

// TableName sets the insert table name for MasterKategoriTransaksi
func (MasterKategoriTransaksi) TableName() string {
	return "master_kategori_transaksis"
}
