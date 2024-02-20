package repository

import "api-hotel-booking/internal/models"

type MerchantAccountRepository interface {
	Save(models.MerchantAccount) (models.MerchantAccount, error)
	FindAll(models.Pagination) (models.Pagination, error)
	FindByID(uint32) (models.MerchantAccount, error)
	Update(uint32, models.MerchantAccount) (int64, error)
	Delete(uint32) (int64, error)
}
