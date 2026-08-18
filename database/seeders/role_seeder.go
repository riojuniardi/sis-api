package seeders

import (
	"sis-api/models"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	roles := []models.Role{
		{Name: "admin"},
		{Name: "guru"},
		{Name: "siswa"},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, models.Role{Name: role.Name}).Error; err != nil {
			return err
		}
	}
	return nil
}
