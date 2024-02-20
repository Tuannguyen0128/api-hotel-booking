package models

type MerchantAccount struct {
	ID           uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Description  string `gorm:"size:100" json:"description"`
	MerchantCode string `gorm:"size:10;not null;unique" json:"merchant_code"`
	//Teammembers  []TeamMember `gorm:"foreinkey:MerchantCode;references:MerchantCode" json:"teammembers,omitempty"`
}
