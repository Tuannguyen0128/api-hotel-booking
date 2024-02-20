package database

import (
	"ProjectPractice/src/config"

	//"github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.DBURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, err
}