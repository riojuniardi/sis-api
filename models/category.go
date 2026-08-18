package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `json:"name"`
	Code string `json:"code" gorm:"unique"`
}

type CategoryInput struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}
