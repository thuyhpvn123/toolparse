package database

import (
	"log"

	"github.com/meta-node-blockchain/mail/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitializeDB(connectionString string) {
	var err error
	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.File{},
		&models.Email{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Failed to close database connection: %v", err)
	}
	sqlDB.Close()
}