package paginator

import (
	"ProjectPractice/src/api/models"
	"fmt"
	"math"

	"gorm.io/gorm"
)

func Paginate(value interface{}, pagination *models.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)
	fmt.Printf("Total rows: %d", pagination.TotalRows)
	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages
	fmt.Printf("Total page: %d", pagination.TotalPages)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}
