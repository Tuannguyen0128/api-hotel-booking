package auto

import (
	"api-hotel-booking/internal/database"
	"api-hotel-booking/internal/models"
	"api-hotel-booking/internal/utils/console"

	"log"
)

func Load() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)

	}

	err = db.Debug().AutoMigrate(&models.Account{})
	if err != nil {
		log.Fatal(err)
	}

	//err =db.Debug().Model(&models.TeamMember)
	for _, account := range accounts {
		err = db.Debug().Model(&models.Account{}).Create(&account).Error
		if err != nil {
			log.Fatal(err)
		}
		console.Pretty(account)
	}
}
