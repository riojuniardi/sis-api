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

type UserRegisterInput struct {
	Name     string `json:"name" binding:"required"`
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   int    `json:"role_id" binding:"required"`
}

type UserLoginInput struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}
