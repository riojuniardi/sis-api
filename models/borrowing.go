package models

import (
	"time"

	"gorm.io/gorm"
)

type BorrowStatus string

const (
	StatusPending        BorrowStatus = "pending"
	StatusDisetujui      BorrowStatus = "disetujui"
	StatusDitolak        BorrowStatus = "ditolak"
	BorrowStatusDipinjam BorrowStatus = BorrowStatus(StatusDipinjam) // Mengambil nilai "dipinjam" dari StatusDipinjam milik ItemStatus
	StatusDikembalikan   BorrowStatus = "dikembalikan"
	StatusTerlambat      BorrowStatus = "terlambat"
)

type Borrowing struct {
	gorm.Model

	ItemID          uint           `json:"item_id" gorm:"not null;index"`
	UserID          uint           `json:"user_id" gorm:"not null;index"`
	BorrowDate      time.Time      `json:"borrow_date" gorm:"not null"`
	DueDate         time.Time      `json:"due_date" gorm:"not null"`
	ReturnDate      *time.Time     `json:"return_date,omitempty"`
	ConditionBefore ItemCondition  `json:"condition_before" gorm:"type:varchar(20);default:'baik';not null"`
	ConditionAfter  *ItemCondition `json:"condition_after,omitempty" gorm:"type:varchar(20)"`
	Status          BorrowStatus   `json:"status" gorm:"type:varchar(20);default:'pending';not null"`
	Notes           *string        `json:"notes,omitempty" gorm:"type:text"`

	// Relasi Database
	Item Item `json:"item,omitempty" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	User User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type BorrowingInput struct {
	ItemID          uint           `json:"item_id" binding:"required"`
	UserID          uint           `json:"user_id" binding:"required"`
	BorrowDate      string         `json:"borrow_date" binding:"required"`
	DueDate         string         `json:"due_date" binding:"required"`
	ReturnDate      *string        `json:"return_date"`
	ConditionBefore ItemCondition  `json:"condition_before"`
	ConditionAfter  *ItemCondition `json:"condition_after"`
	Status          BorrowStatus   `json:"status"`
	Notes           *string        `json:"notes"`
}

type BorrowingApprovalInput struct {
	Status BorrowStatus `json:"status" binding:"required"`
	Notes  *string      `json:"notes"`
}

type BorrowingReturnInput struct {
	ConditionAfter ItemCondition `json:"condition_after" binding:"required"`
	Notes          *string       `json:"notes"`
}
