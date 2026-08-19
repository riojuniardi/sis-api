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

func RegisterUser(c *gin.Context) {
	var input models.UserRegisterInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	heshedPassword, errHes := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if errHes != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User name mungkin sudah terdaftar",
		})
		return
	}

	config.DB.Preload("Role").First(&user, user.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil registerasi user",
		"user": gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"user_name": user.UserName,
			"role":      user.Role.Name,
		},
	})
}

func LoginUser(c *gin.Context) {
	var input models.UserLoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	err := config.DB.Preload("Role").Where("user_name = ?", input.UserName).First(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User name tidak terdaftar"})
		return
	}

	errMatchPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if errMatchPassword != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
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

func GetCurrentUser(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Tidak auntentication",
		})
		return
	}

	var user models.User

	userData := config.DB.Preload("Role").First(&user, userId).Error
	if userData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":   user.ID,
			"name": user.Name,
			"role": user.Role.Name,
		},
	})
}

func GetAllUsers(c *gin.Context) {
	var users []models.User
	result := config.DB.Preload("Role").Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data user",
		"users":   userResponses,
	})
}

func GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	result := config.DB.Preload("Role").First(&user, userID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data user",
		"user": gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"user_name": user.UserName,
			"role":      user.Role.Name,
		},
	})
}
