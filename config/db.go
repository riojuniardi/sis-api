package config

import (
	"log"
	"os"
	"sis-api/database/seeders"
	"sis-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Environment variabel belum di isi")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Gagal terkoneksi di database", err)
	}

	err = database.AutoMigrate(&models.Role{}, &models.User{}, &models.Category{}, &models.Location{}, &models.Item{})
	if err != nil {
		log.Fatal("Gagal melakukan migration database :", err)
	}

	DB = database

	if err := seeders.SeedUsers(DB); err != nil {
		log.Fatal("Gagal menjalankan seeder:", err)
	}
	log.Println("Berhasil terkoneksi ke Database")
}
