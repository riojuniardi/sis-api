package controllers

import (
	"errors"
	"net/http"
	"sis-api/config"
	"sis-api/models"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CREATE AUDIT (Mencatat Audit Kondisi Barang Baru)
func CreateConditionAudit(c *gin.Context) {
	var input models.ConditionAuditInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Ambil data item terkini untuk mendapatkan ConditionBefore
	var item models.Item
	if err := config.DB.WithContext(c.Request.Context()).First(&item, input.ItemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Item ID tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengecek data item"})
		return
	}

	// Standardize reason jika kosong
	if input.Reason == "" {
		input.Reason = models.ReasonRoutineCheck
	}

	conditionBefore := item.Condition
	var audit models.ConditionAudit

	// 2. Transaksi Database: Buat Log Audit & Update Kondisi Item
	err := config.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		audit = models.ConditionAudit{
			ItemID:          input.ItemID,
			UserID:          input.UserID,
			BorrowingID:     input.BorrowingID,
			ConditionBefore: conditionBefore,
			ConditionAfter:  input.ConditionAfter,
			Reason:          input.Reason,
			Notes:           input.Notes,
		}

		if err := tx.Create(&audit).Error; err != nil {
			return err
		}

		// Update kondisi barang di tabel items
		updates := map[string]interface{}{
			"condition": input.ConditionAfter,
		}

		// Jika barang dinilai 'rusak_berat', otomatis ubah status ketersediaan item ke 'perbaikan'
		if input.ConditionAfter == models.ConditionRusakBerat {
			updates["status"] = models.StatusPerbaikan
		} else if conditionBefore == models.ConditionRusakBerat && input.ConditionAfter == models.ConditionBaik && item.Status == models.StatusPerbaikan {
			// Jika barang selesai diperbaiki (dari rusak berat ke baik), ubah status ketersediaan kembali ke 'tersedia'
			updates["status"] = models.StatusTersedia
		}

		if err := tx.Model(&models.Item{}).Where("id = ?", input.ItemID).Updates(updates).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") || strings.Contains(err.Error(), "23503") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID atau Borrowing ID tidak valid"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan audit kondisi barang"})
		return
	}

	config.DB.WithContext(c.Request.Context()).
		Preload("Item").
		Preload("User").
		Preload("Borrowing").
		First(&audit, audit.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Audit kondisi barang berhasil dicatat",
		"audit":   audit,
	})
}

// READ ALL AUDITS (Bisa Filter Query Params ?item_id=1)
func GetAllConditionAudits(c *gin.Context) {
	var audits []models.ConditionAudit
	db := config.DB.WithContext(c.Request.Context()).Preload("Item").Preload("User").Preload("Borrowing")

	itemID := c.Query("item_id")
	if itemID != "" {
		db = db.Where("item_id = ?", itemID)
	}

	userID := c.Query("user_id")
	if userID != "" {
		db = db.Where("user_id = ?", userID)
	}

	if err := db.Order("created_at desc").Find(&audits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data audit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"audits": audits,
	})
}

// READ DETAIL AUDIT BY ID
func GetConditionAuditByID(c *gin.Context) {
	id := c.Param("id")
	var audit models.ConditionAudit

	if err := config.DB.WithContext(c.Request.Context()).
		Preload("Item").
		Preload("User").
		Preload("Borrowing").
		First(&audit, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data audit tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data audit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"audit": audit,
	})
}

// READ HISTORI AUDIT UNTUK SPESIFIK ITEM
func GetAuditsByItemID(c *gin.Context) {
	itemID := c.Param("item_id")
	var audits []models.ConditionAudit

	if err := config.DB.WithContext(c.Request.Context()).
		Where("item_id = ?", itemID).
		Preload("User").
		Preload("Borrowing").
		Order("created_at desc").
		Find(&audits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat audit item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"item_id": itemID,
		"audits":  audits,
	})
}

// DELETE AUDIT LOG
func DeleteConditionAudit(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.WithContext(c.Request.Context()).Delete(&models.ConditionAudit{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus log audit"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data audit tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil menghapus log audit",
	})
}
