package controllers

import (
	"net/http"
	"os"
	"sis-api/config"
	"sis-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthInputRegister struct {
	Name     string `json:"name" binding:"required"`
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   int    `json:"role_id" binding:"required"`
}

type AuthInputLogin struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterUser(context *gin.Context) {
	var input AuthInputRegister

	err := context.ShouldBindJSON(&input)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	heshedPassword, errHes := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if errHes != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal encripsi Password",
		})
		return
	}

	user := models.User{
		Name:     input.Name,
		UserName: input.UserName,
		Password: string(heshedPassword),
		RoleID:   input.RoleID,
	}

	userCreate := config.DB.Create(&user).Error
	if userCreate != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "User name mungkin sudah terdaftar",
		})
		return
	}

	// Load role relationship
	config.DB.Preload("Role").First(&user, user.ID)

	context.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil registerasi user",
		"user": gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"user_name": user.UserName,
			"role":      user.Role.Name,
		},
	})
}

func LoginUser(context *gin.Context) {
	var input AuthInputLogin

	err := context.ShouldBindJSON(&input)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User
	userData := config.DB.Where("user_name = ?", input.UserName).First(&user).Error
	if userData != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "User name tidak terdaftar",
		})
		return
	}

	errMatchPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if errMatchPassword != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "Password salah",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role.Name,
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat token",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"Message": "Login berhasil",
		"token":   tokenString,
		"user": gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"user_name": user.UserName,
			"role":      user.Role.Name,
		},
	})
}

func GetCurrentUser(context *gin.Context) {
	userId, exists := context.Get("userID")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "Tidak auntentication",
		})
		return
	}

	var user models.User

	userData := config.DB.Preload("Role").First(&user, userId).Error
	if userData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "User tidak ditemukan",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":   user.ID,
			"name": user.Name,
			"role": user.Role.Name,
		},
	})
}

// GetAllUsers - Mendapatkan semua user (Admin only)
func GetAllUsers(context *gin.Context) {
	var users []models.User
	result := config.DB.Preload("Role").Find(&users)

	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data user",
		})
		return
	}

	var userResponses []gin.H
	for _, user := range users {
		userResponses = append(userResponses, gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"user_name": user.UserName,
			"role":      user.Role.Name,
		})
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data user",
		"users":   userResponses,
	})
}

// GetUserByID - Mendapatkan user berdasarkan ID (Admin only)
func GetUserByID(context *gin.Context) {
	userID := context.Param("id")

	var user models.User
	result := config.DB.Preload("Role").First(&user, userID)

	if result.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "User tidak ditemukan",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data user",
		"user": gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"user_name": user.UserName,
			"role":      user.Role.Name,
		},
	})
}
