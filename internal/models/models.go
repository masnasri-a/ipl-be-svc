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
	ID          uint       `json:"id" gorm:"primarykey"`
	DocumentID  string     `json:"document_id" gorm:"column:document_id"`
	NamaMenu    string     `json:"nama_menu" gorm:"column:nama_menu"`
	KodeMenu    string     `json:"kode_menu" gorm:"column:kode_menu"`
	UrutanMenu  *int       `json:"urutan_menu" gorm:"column:urutan_menu"`
	IsActive    *bool      `json:"is_active" gorm:"column:is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
	Locale      *string    `json:"locale"`
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

// UserDetail represents the user detail response
type UserDetail struct {
	ID           uint   `json:"id" gorm:"column:id"`
	NamaPenghuni string `json:"nama_penghuni" gorm:"column:nama_penghuni"`
	NoHP         string `json:"no_hp" gorm:"column:no_hp"`
	NoTelp       string `json:"no_telp" gorm:"column:no_telp"`
	DocumentID   string `json:"document_id" gorm:"column:document_id"`
	Email        string `json:"email" gorm:"column:email"`
	UserID       uint   `json:"user_id" gorm:"column:user_id"`
	RoleName     string `json:"role_name" gorm:"column:name"`
	RoleID       uint   `json:"role_id" gorm:"column:role_id"`
	RoleType     string `json:"role_type" gorm:"column:role_type"`
}

// TableName sets the insert table name for UserDetail
func (UserDetail) TableName() string {
	return "profiles"
}

// RoleMenu represents the role_menus table
type RoleMenu struct {
	ID          uint       `json:"id" gorm:"primarykey"`
	DocumentID  *string    `json:"document_id" gorm:"column:document_id"`
	RoleMenuOrd *float64   `json:"role_menu_ord" gorm:"column:role_menu_ord"`
	IsActive    *bool      `json:"is_active" gorm:"column:is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedByID *int       `json:"created_by_id"`
	UpdatedByID *int       `json:"updated_by_id"`
	// Relationships
	MasterMenus []MasterMenu `json:"master_menus,omitempty" gorm:"many2many:role_menus_master_menu_lnk;joinForeignKey:role_menu_id;joinReferences:master_menu_id"`
	Roles       []Role       `json:"roles,omitempty" gorm:"many2many:role_menus_role_lnk;joinForeignKey:role_menu_id;joinReferences:role_id"`
}

// TableName sets the insert table name for RoleMenu
func (RoleMenu) TableName() string {
	return "role_menus"
}

// Role represents the up_roles table
type Role struct {
	ID          uint       `json:"id" gorm:"primarykey"`
	DocumentID  *string    `json:"document_id" gorm:"column:document_id"`
	Name        *string    `json:"name" gorm:"column:name"`
	Description *string    `json:"description" gorm:"column:description"`
	Type        *string    `json:"type" gorm:"column:type"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedByID *int       `json:"created_by_id"`
	UpdatedByID *int       `json:"updated_by_id"`
}

// TableName sets the insert table name for Role
func (Role) TableName() string {
	return "up_roles"
}

// RoleMenuMasterMenuLink represents the role_menus_master_menu_lnk table
type RoleMenuMasterMenuLink struct {
	ID           uint     `json:"id" gorm:"primarykey"`
	RoleMenuID   uint     `json:"role_menu_id" gorm:"column:role_menu_id"`
	MasterMenuID uint     `json:"master_menu_id" gorm:"column:master_menu_id"`
	RoleMenuOrd  *float64 `json:"role_menu_ord" gorm:"column:role_menu_ord"`
}

// TableName sets the insert table name for RoleMenuMasterMenuLink
func (RoleMenuMasterMenuLink) TableName() string {
	return "role_menus_master_menu_lnk"
}

// RoleMenuRoleLink represents the role_menus_role_lnk table
type RoleMenuRoleLink struct {
	ID          uint     `json:"id" gorm:"primarykey"`
	RoleMenuID  uint     `json:"role_menu_id" gorm:"column:role_menu_id"`
	RoleID      uint     `json:"role_id" gorm:"column:role_id"`
	RoleMenuOrd *float64 `json:"role_menu_ord" gorm:"column:role_menu_ord"`
}

// TableName sets the insert table name for RoleMenuRoleLink
func (RoleMenuRoleLink) TableName() string {
	return "role_menus_role_lnk"
}
