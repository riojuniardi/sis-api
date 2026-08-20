package models

import (
	"time"

	"gorm.io/gorm"
)

type ItemCondition string

const (
	ConditionBaik        ItemCondition = "baik"
	ConditionRusakRingan ItemCondition = "rusak_ringan"
	ConditionRusakBerat  ItemCondition = "rusak_berat"
)

type ItemStatus string

const (
	StatusTersedia  ItemStatus = "tersedia"
	StatusDipinjam  ItemStatus = "dipinjam"
	StatusPerbaikan ItemStatus = "perbaikan"
	StatusAfkir     ItemStatus = "afkir"
)

type Item struct {
	gorm.Model

	ItemCode      string        `json:"item_code" gorm:"uniqueIndex"`
	CategoryID    uint          `json:"category_id"`
	LocationID    uint          `json:"location_id""`
	Name          string        `json:"name"`
	SourceOfFunds string        `json:"source_of_funds"`
	PurchaseDate  time.Time     `json:"purchase_date" gorm:"type:date;not null"`
	PurchasePrice float64       `json:"purchase_price" gorm:"type:decimal(20,2);not null;default:0"`
	Condition     ItemCondition `json:"condition" gorm:"default:'baik';not null"`
	Status        ItemStatus    `json:"status" gorm:"default:'tersedia';not null"`

	Category Category `json:"category,omitempty" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Location Location `json:"location,omitempty" gorm:"foreignKey:LocationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type ItemInput struct {
	ItemCode      string        `json:"item_code" binding:"required"`
	CategoryID    uint          `json:"category_id" binding:"required"`
	LocationID    uint          `json:"location_id" binding:"required"`
	Name          string        `json:"name" binding:"required"`
	SourceOfFunds string        `json:"source_of_funds" binding:"required"`
	PurchaseDate  string        `json:"purchase_date" binding:"required"`
	PurchasePrice float64       `json:"purchase_price"`
	Condition     ItemCondition `json:"condition"`
	Status        ItemStatus    `json:"status"`
}
