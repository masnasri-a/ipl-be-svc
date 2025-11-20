package models

// BillingKategoriTransaksiLink represents the billings_master_kategori_transaksi_lnk table
type BillingKategoriTransaksiLink struct {
	ID                        uint `json:"id" gorm:"primarykey"`
	BillingID                 uint `json:"t_billing_id" gorm:"column:t_billing_id"`
	MasterKategoriTransaksiID uint `json:"master_kategori_transaksi_id" gorm:"column:master_kategori_transaksi_id"`
}

// TableName sets the insert table name for BillingKategoriTransaksiLink
func (BillingKategoriTransaksiLink) TableName() string {
	return "billings_master_kategori_transaksi_lnk"
}
