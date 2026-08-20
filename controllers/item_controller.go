package controllers

import (
	"errors"
	"net/http"
	"sis-api/config"
	"sis-api/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateItem(c *gin.Context) {
	var input models.ItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	purchaseDate, err := time.Parse("2006-01-02", input.PurchaseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid. Gunakan YYYY-MM-DD"})
		return
	}

	if input.Condition == "" {
		input.Condition = models.ConditionBaik
	}
	if input.Status == "" {
		input.Status = models.StatusTersedia
	}

	item := models.Item{
		ItemCode:      input.ItemCode,
		CategoryID:    input.CategoryID,
		LocationID:    input.LocationID,
		Name:          input.Name,
		SourceOfFunds: input.SourceOfFunds,
		PurchaseDate:  purchaseDate,
		PurchasePrice: input.PurchasePrice,
		Condition:     input.Condition,
		Status:        input.Status,
	}

	if err := config.DB.WithContext(c.Request.Context()).Create(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "23505") {
			c.JSON(http.StatusConflict, gin.H{"error": "Kode item sudah digunakan"})
			return
		}
		if strings.Contains(err.Error(), "foreign key constraint") || strings.Contains(err.Error(), "23503") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category ID atau Location ID tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan item"})
		return
	}

	config.DB.WithContext(c.Request.Context()).Preload("Category").Preload("Location").First(&item, item.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil menambah item",
		"item":    item,
	})
}

func GetAllItems(c *gin.Context) {
	var items []models.Item

	if err := config.DB.WithContext(c.Request.Context()).Preload("Category").Preload("Location").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
	})
}

func GetItemByID(c *gin.Context) {
	id := c.Param("id")
	var item models.Item

	if err := config.DB.WithContext(c.Request.Context()).Preload("Category").Preload("Location").First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"item": item,
	})
}

func UpdateItem(c *gin.Context) {
	id := c.Param("id")
	var item models.Item

	if err := config.DB.WithContext(c.Request.Context()).First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari item"})
		return
	}

	var input models.ItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	purchaseDate, err := time.Parse("2006-01-02", input.PurchaseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid. Gunakan YYYY-MM-DD"})
		return
	}

	item.ItemCode = input.ItemCode
	item.CategoryID = input.CategoryID
	item.LocationID = input.LocationID
	item.Name = input.Name
	item.SourceOfFunds = input.SourceOfFunds
	item.PurchaseDate = purchaseDate
	item.PurchasePrice = input.PurchasePrice
	item.Condition = input.Condition
	item.Status = input.Status

	if err := config.DB.WithContext(c.Request.Context()).Save(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "23505") {
			c.JSON(http.StatusConflict, gin.H{"error": "Kode item sudah digunakan"})
			return
		}
		if strings.Contains(err.Error(), "foreign key constraint") || strings.Contains(err.Error(), "23503") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category ID atau Location ID tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui item"})
		return
	}

	config.DB.WithContext(c.Request.Context()).Preload("Category").Preload("Location").First(&item, item.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil memperbarui item",
		"item":    item,
	})
}

func DeleteItem(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.WithContext(c.Request.Context()).Delete(&models.Item{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus item"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil menghapus item",
	})
}
