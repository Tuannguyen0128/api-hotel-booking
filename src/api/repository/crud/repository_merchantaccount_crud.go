package crud

import (
	"ProjectPractice/src/api/models"
	"ProjectPractice/src/api/utils/channels"
	"ProjectPractice/src/api/utils/paginator"
	"errors"

	"gorm.io/gorm"
)

type repositoryMerchantAccountCRUD struct {
	db *gorm.DB
}

func NewRepositoryMerchantAccountCRUD(db *gorm.DB) *repositoryMerchantAccountCRUD {
	return &repositoryMerchantAccountCRUD{db}
}
func (r *repositoryMerchantAccountCRUD) Save(merchantaccount models.MerchantAccount) (models.MerchantAccount, error) {
	var err error
	done := make(chan bool)
	go func(ch chan<- bool) {
		err = r.db.Debug().Model(&models.MerchantAccount{}).Create(&merchantaccount).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)
	if channels.OK(done) {
		return merchantaccount, nil
	}
	return models.MerchantAccount{}, err
}
func (r *repositoryMerchantAccountCRUD) FindAll(pagination models.Pagination) (models.Pagination, error) {
	var err error
	merchantaccount := []models.MerchantAccount{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		err = r.db.Debug().Model(&models.MerchantAccount{}).Scopes(paginator.Paginate(merchantaccount, &pagination, r.db)).Find(&merchantaccount).Error
		pagination.Rows = merchantaccount
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)
	if channels.OK(done) {
		return pagination, nil
	}
	return pagination, err
}
func (r *repositoryMerchantAccountCRUD) Update(uid uint32, merchantaccount models.MerchantAccount) (int64, error) {
	var rs *gorm.DB

	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		rs = r.db.Debug().Model(&models.MerchantAccount{}).Where("id=?", uid).Take(&models.MerchantAccount{}).UpdateColumns(
			map[string]interface{}{
				"description":   merchantaccount.Description,
				"merchant_code": merchantaccount.MerchantCode,
			},
		)

		ch <- true
	}(done)
	if channels.OK(done) {
		if rs.Error != nil {
			return 0, rs.Error
		}
		return rs.RowsAffected, nil

	}
	return 0, rs.Error
}
func (r *repositoryMerchantAccountCRUD) Delete(uid uint32) (int64, error) {
	var rs *gorm.DB

	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		rs = r.db.Debug().Model(&models.MerchantAccount{}).Where("id=?", uid).Delete(&models.MerchantAccount{})

		ch <- true
	}(done)
	if channels.OK(done) {
		if rs.Error != nil {
			return 0, rs.Error
		}
		return rs.RowsAffected, nil

	}
	return 0, rs.Error
}
func (r *repositoryMerchantAccountCRUD) FindByID(uid uint32) (models.MerchantAccount, error) {

	var err error
	merchant := models.MerchantAccount{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		err = r.db.Debug().Model(&models.MerchantAccount{}).Where("id = ?", uid).Find(&merchant).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)
	errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil {
		return models.MerchantAccount{}, errors.New("record not found")
	}
	if channels.OK(done) {
		return merchant, nil
	}
	return merchant, err
}
