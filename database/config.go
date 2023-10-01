package database

import (
	"log"
	"os"

	"github.com/crocox/final-project/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectionDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env file", err)
	}

	dsn := os.Getenv("DB_URL")
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Could not load the database", err)
	}

	DB.AutoMigrate(
		&models.User{},
		&models.Photo{},
	)

	log.Println("Database Connected")

	return DB
}
