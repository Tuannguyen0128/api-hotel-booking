package repository

import (
	"api-hotel-booking/internal/models"
	"api-hotel-booking/internal/utils/channels"
	"api-hotel-booking/internal/utils/paginator"
	"errors"
	"time"

	"gorm.io/gorm"
)

type AccountRepo interface {
	Save(models.Account) (models.Account, error)
	FindAll(models.Pagination) (models.Pagination, error)
	FindByEmail(string) (models.Account, error)
	Update(uint32, models.Account) (int64, error)
	Delete(uint32) (int64, error)
	FindByMerchantCode(string) ([]models.Account, error)
}

type repositoryAccountCRUD struct {
	db *gorm.DB
}

func NewRepositoryAccountCRUD(db *gorm.DB) *repositoryAccountCRUD {
	return &repositoryAccountCRUD{db}
}
func (r *repositoryAccountCRUD) Save(account models.Account) (models.Account, error) {
	var err error
	done := make(chan bool)
	go func(ch chan<- bool) {
		err = r.db.Debug().Model(&models.Account{}).Create(&account).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)
	if channels.OK(done) {
		return account, nil
	}
	return models.Account{}, err
}
func (r *repositoryAccountCRUD) FindAll(pagination models.Pagination) (models.Pagination, error) {
	var err error
	accounts := []models.Account{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		err = r.db.Debug().Model(&models.Account{}).Scopes(paginator.Paginate(accounts, &pagination, r.db)).Find(&accounts).Error
		pagination.Rows = accounts
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
func (r *repositoryAccountCRUD) FindByEmail(email string) (models.Account, error) {
	var err error
	account := models.Account{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		err = r.db.Debug().Model(&models.Account{}).Where("email = ?", email).Take(&account).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)
	errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil {
		return models.Account{}, errors.New("email not found")
	}
	if channels.OK(done) {
		return account, nil
	}
	return account, err
}
func (r *repositoryAccountCRUD) Update(uid uint32, account models.Account) (int64, error) {
	var rs *gorm.DB

	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		rs = r.db.Debug().Model(&models.Account{}).Where("id=?", uid).Take(&models.Account{}).UpdateColumns(
			map[string]interface{}{
				"fullname":   account.Fullname,
				"email":      account.Email,
				"updated_at": time.Now(),
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
func (r *repositoryAccountCRUD) Delete(uid uint32) (int64, error) {
	var rs *gorm.DB

	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		rs = r.db.Debug().Model(&models.Account{}).Where("id=?", uid).Delete(&models.Account{})

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
func (r *repositoryAccountCRUD) FindByMerchantCode(merchantcode string) ([]models.Account, error) {
	var err error
	teammembers := []models.Account{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		err = r.db.Debug().Model(&models.Account{}).Where("merchant_code = ?", merchantcode).Find(&teammembers).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)
	errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil {
		return []models.Account{}, errors.New("email not found")
	}
	if channels.OK(done) {
		return teammembers, nil
	}
	return teammembers, err
}
