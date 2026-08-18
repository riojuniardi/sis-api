package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	RoleID   int    `json:"role_id"`
	Role     Role   `json:"-" gorm:"foreignKey:RoleID"`
}
