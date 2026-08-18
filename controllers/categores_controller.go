package controllers

import (
	"net/http"
	"sis-api/config"
	"sis-api/models"
	"strings"

	"github.com/gin-gonic/gin"
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

	if err := config.DB.Create(&category).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "23505") {
			c.JSON(http.StatusConflict, gin.H{"error": "Kode kategori sudah digunakan"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan kategori"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": category})
}

func GetAllCategories(c *gin.Context) {
	var categories []models.Category
	config.DB.Find(&categories)

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func GetCategoryByID(c *gin.Context) {
	var category models.Category
	id := c.Param("id")

	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

func UpdateCategory(c *gin.Context) {
	var category models.Category
	id := c.Param("id")

	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	var input models.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Update data
	config.DB.Model(&category).Updates(models.Category{
		Name: input.Name,
		Code: input.Code,
	})

	c.JSON(http.StatusOK, gin.H{"data": category})
}

func DeleteCategory(c *gin.Context) {
	var category models.Category
	id := c.Param("id")

	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	config.DB.Delete(&category)

	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus"})
}
