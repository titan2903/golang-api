package config

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/transaction"
	"bwastartup/user"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	// dsnMaster := fmt.Sprintf(
	// 	"host=%s user=%s password=%s dbname=%s port=%s",
	// 	helper.GoDotEnvVariable("DB_HOST"), helper.GoDotEnvVariable("DB_USER"), helper.GoDotEnvVariable("DB_PASSWORD"), helper.GoDotEnvVariable("DB_NAME"), helper.GoDotEnvVariable("DB_PORT"),
	// )
	dsnMaster := helper.GoDotEnvVariable("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsnMaster), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	if err = db.AutoMigrate(
		&user.User{},
		&campaign.Campaign{},
		&transaction.Transaction{},
	); err != nil {
		log.Fatal(err)
	}

	log.Println("success connect to database")

	return db
}
