package auto

import (
	"ProjectPractice/src/api/database"
	"ProjectPractice/src/api/models"
	"ProjectPractice/src/api/utils/console"

	"log"
)

func Load() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)

	}

	err = db.Debug().AutoMigrate(&models.MerchantAccount{}, &models.TeamMember{})
	if err != nil {
		log.Fatal(err)
	}
	for _, merchant := range merchantaccount {
		err = db.Debug().Model(&models.MerchantAccount{}).Create(&merchant).Error
		if err != nil {
			log.Fatal(err)
		}
		console.Pretty(merchant)
	}
	//err =db.Debug().Model(&models.TeamMember)
	for _, teammember := range teammembers {
		err = db.Debug().Model(&models.TeamMember{}).Create(&teammember).Error
		if err != nil {
			log.Fatal(err)
		}
		console.Pretty(teammember)
	}
}
