package models

import (
	"gorm.io/gorm"
)

type AuditReason string

const (
	ReasonRoutineCheck AuditReason = "pemeriksaan_rutin"
	ReasonBorrowReturn AuditReason = "pengembalian_pinjaman"
	ReasonDamageReport AuditReason = "laporan_kerusakan"
	ReasonPostRepair   AuditReason = "selesai_perbaikan"
)

type ConditionAudit struct {
	gorm.Model

	ItemID          uint          `json:"item_id" gorm:"not null;index"`
	UserID          uint          `json:"user_id" gorm:"not null;index"`       // Auditor / Petugas yang memeriksa
	BorrowingID     *uint         `json:"borrowing_id,omitempty" gorm:"index"` // Opsional, jika audit terkait peminjaman
	ConditionBefore ItemCondition `json:"condition_before" gorm:"type:varchar(20);not null"`
	ConditionAfter  ItemCondition `json:"condition_after" gorm:"type:varchar(20);not null"`
	Reason          AuditReason   `json:"reason" gorm:"type:varchar(30);default:'pemeriksaan_rutin';not null"`
	Notes           *string       `json:"notes,omitempty" gorm:"type:text"`

	// Relasi Database
	Item      Item       `json:"item,omitempty" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	User      User       `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Borrowing *Borrowing `json:"borrowing,omitempty" gorm:"foreignKey:BorrowingID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ConditionAuditInput struct {
	ItemID         uint          `json:"item_id" binding:"required"`
	UserID         uint          `json:"user_id" binding:"required"`
	BorrowingID    *uint         `json:"borrowing_id"`
	ConditionAfter ItemCondition `json:"condition_after" binding:"required"` // baik / rusak_ringan / rusak_berat
	Reason         AuditReason   `json:"reason"`                             // pemeriksaan_rutin / laporan_kerusakan / selesai_perbaikan / pengembalian_pinjaman
	Notes          *string       `json:"notes"`
}
