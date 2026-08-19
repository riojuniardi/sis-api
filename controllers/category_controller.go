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

func CreateCategory(c *gin.Context) {
	var input models.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
		Code: input.Code,
		Name: input.Name,
	}

	if err := config.DB.WithContext(c.Request.Context()).Create(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "23505") {
			c.JSON(http.StatusConflict, gin.H{"error": "Kode kategori sudah digunakan"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan kategori"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Berhasil menambah kategori",
		"category": category,
	})
}

func GetAllCategories(c *gin.Context) {
	var categories []models.Category

	if err := config.DB.WithContext(c.Request.Context()).Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kategori"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

func GetCategoryByID(c *gin.Context) {
	var category models.Category
	id := c.Param("id")

	if err := config.DB.WithContext(c.Request.Context()).First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kategori"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": category,
	})
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := config.DB.WithContext(c.Request.Context()).First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari kategori"})
		return
	}

	var input models.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.Code = input.Code
	category.Name = input.Name

	if err := config.DB.WithContext(c.Request.Context()).Save(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "23505") {
			c.JSON(http.StatusConflict, gin.H{"error": "Kode kategori sudah digunakan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui kategori"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Berhasil memperbarui kategori",
		"category": category,
	})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.WithContext(c.Request.Context()).Delete(&models.Category{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil menghapus kategori",
	})
}
