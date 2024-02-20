package crud

import (
	"api-hotel-booking/internal/models"
	"api-hotel-booking/internal/utils/channels"
	"api-hotel-booking/internal/utils/paginator"
	"errors"
	"time"

	"gorm.io/gorm"
)

type repositoryTeamMemberCRUD struct {
	db *gorm.DB
}

func NewRepositoryTeamMemberCRUD(db *gorm.DB) *repositoryTeamMemberCRUD {
	return &repositoryTeamMemberCRUD{db}
}
func (r *repositoryTeamMemberCRUD) Save(teammember models.TeamMember) (models.TeamMember, error) {
	var err error
	done := make(chan bool)
	go func(ch chan<- bool) {
		err = r.db.Debug().Model(&models.TeamMember{}).Create(&teammember).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)
	if channels.OK(done) {
		return teammember, nil
	}
	return models.TeamMember{}, err
}
func (r *repositoryTeamMemberCRUD) FindAll(pagination models.Pagination) (models.Pagination, error) {
	var err error
	teammembers := []models.TeamMember{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		err = r.db.Debug().Model(&models.TeamMember{}).Scopes(paginator.Paginate(teammembers, &pagination, r.db)).Find(&teammembers).Error
		pagination.Rows = teammembers
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
func (r *repositoryTeamMemberCRUD) FindByEmail(email string) (models.TeamMember, error) {
	var err error
	teammember := models.TeamMember{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		err = r.db.Debug().Model(&models.TeamMember{}).Where("email = ?", email).Take(&teammember).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)
	errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil {
		return models.TeamMember{}, errors.New("email not found")
	}
	if channels.OK(done) {
		return teammember, nil
	}
	return teammember, err
}
func (r *repositoryTeamMemberCRUD) Update(uid uint32, teammember models.TeamMember) (int64, error) {
	var rs *gorm.DB

	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		rs = r.db.Debug().Model(&models.TeamMember{}).Where("id=?", uid).Take(&models.TeamMember{}).UpdateColumns(
			map[string]interface{}{
				"fullname":   teammember.Fullname,
				"email":      teammember.Email,
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
func (r *repositoryTeamMemberCRUD) Delete(uid uint32) (int64, error) {
	var rs *gorm.DB

	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		rs = r.db.Debug().Model(&models.TeamMember{}).Where("id=?", uid).Delete(&models.TeamMember{})

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
func (r *repositoryTeamMemberCRUD) FindByMerchantCode(merchantcode string) ([]models.TeamMember, error) {
	var err error
	teammembers := []models.TeamMember{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		err = r.db.Debug().Model(&models.TeamMember{}).Where("merchant_code = ?", merchantcode).Find(&teammembers).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)
	errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil {
		return []models.TeamMember{}, errors.New("email not found")
	}
	if channels.OK(done) {
		return teammembers, nil
	}
	return teammembers, err
}
