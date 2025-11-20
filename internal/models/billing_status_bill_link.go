package models

// BillingStatusBillLink represents the billings_status_bill_lnk table
type BillingStatusBillLink struct {
	ID                    uint `json:"id" gorm:"primarykey"`
	BillingID             uint `json:"t_billing_id" gorm:"column:t_billing_id"`
	MasterGeneralStatusID uint `json:"master_general_status_id" gorm:"column:master_general_status_id"`
}

// TableName sets the insert table name for BillingStatusBillLink
func (BillingStatusBillLink) TableName() string {
	return "billings_status_bill_lnk"
}
