package models

import "gorm.io/gorm"

type Location struct {
	gorm.Model
	Name      string `json:"name"`
	Code      string `json:"code" gorm:"unique"`
	Building  string `json:"building"`
	Floor     int    `json:"floor"`
	PicUserID int    `json:"pic_user_id"`
	User      User   `json:"-" gorm:"foreignKey:PicUserID"`
}

type LocationInput struct {
	Name      string `json:"name" binding:"required"`
	Code      string `json:"code" binding:"required"`
	Building  string `json:"building" binding:"required"`
	Floor     int    `json:"floor" binding:"required"`
	PicUserID int    `json:"pic_user_id" binding:"required"`
}
