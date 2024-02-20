package repository

import "api-hotel-booking/internal/models"

type TeamMemberRepository interface {
	Save(models.TeamMember) (models.TeamMember, error)
	FindAll(models.Pagination) (models.Pagination, error)
	FindByEmail(string) (models.TeamMember, error)
	Update(uint32, models.TeamMember) (int64, error)
	Delete(uint32) (int64, error)
	FindByMerchantCode(string) ([]models.TeamMember, error)
}
