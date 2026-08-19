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

func formatLocationResponse(loc models.Location) gin.H {
	return gin.H{
		"id":       loc.ID,
		"name":     loc.Name,
		"code":     loc.Code,
		"building": loc.Building,
		"floor":    loc.Floor,
		"pic_user": loc.User.Name,
	}
}

func CreateLocation(c *gin.Context) {
	var input models.LocationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location := models.Location{
		Name:      input.Name,
		Code:      input.Code,
		Building:  input.Building,
		Floor:     input.Floor,
		PicUserID: input.PicUserID,
	}

	if err := config.DB.WithContext(c.Request.Context()).Create(&location).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "23505") {
			c.JSON(http.StatusConflict, gin.H{"error": "Kode location sudah digunakan"})
			return
		}
		if strings.Contains(err.Error(), "foreign key constraint") || strings.Contains(err.Error(), "23503") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User PIC tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan location"})
		return
	}

	config.DB.WithContext(c.Request.Context()).Preload("User").First(&location, location.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Berhasil menambah location",
		"location": formatLocationResponse(location),
	})
}

func GetLocations(c *gin.Context) {
	var locations []models.Location

	if err := config.DB.WithContext(c.Request.Context()).Preload("User").Find(&locations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data location"})
		return
	}

	var responseList []gin.H
	for _, loc := range locations {
		responseList = append(responseList, formatLocationResponse(loc))
	}

	c.JSON(http.StatusOK, gin.H{
		"locations": responseList,
	})
}

func GetLocationByID(c *gin.Context) {
	id := c.Param("id")
	var location models.Location

	if err := config.DB.WithContext(c.Request.Context()).Preload("User").First(&location, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Location tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"location": formatLocationResponse(location),
	})
}

func UpdateLocation(c *gin.Context) {
	id := c.Param("id")
	var location models.Location

	// Cek apakah data ada di DB
	if err := config.DB.WithContext(c.Request.Context()).First(&location, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Location tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari location"})
		return
	}

	var input models.LocationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location.Name = input.Name
	location.Code = input.Code
	location.Building = input.Building
	location.Floor = input.Floor
	location.PicUserID = input.PicUserID

	if err := config.DB.WithContext(c.Request.Context()).Save(&location).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "23505") {
			c.JSON(http.StatusConflict, gin.H{"error": "Kode location sudah digunakan"})
			return
		}
		if strings.Contains(err.Error(), "foreign key constraint") || strings.Contains(err.Error(), "23503") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User PIC tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui location"})
		return
	}

	config.DB.WithContext(c.Request.Context()).Preload("User").First(&location, location.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Berhasil memperbarui location",
		"location": formatLocationResponse(location),
	})
}

func DeleteLocation(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.WithContext(c.Request.Context()).Delete(&models.Location{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus location"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Location tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil menghapus location",
	})
}
