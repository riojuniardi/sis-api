package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequiredAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			context.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{
					"error": "Akses ditolak, token tidak ditemukan",
				})
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			context.Set("userID", int(claims["sub"].(float64)))
			context.Next()
		} else {
			context.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{
					"error": "Token tidak valid",
				})
		}
	}

}

// RequiredAdmin middleware untuk validasi token dan role admin dari JWT
func RequiredAdmin() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			context.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{
					"error": "Akses ditolak, token tidak ditemukan",
				})
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := int(claims["sub"].(float64))
			context.Set("userID", userID)

			role, ok := claims["role"].(string)
			if !ok {
				context.AbortWithStatusJSON(http.StatusUnauthorized,
					gin.H{
						"error": "Role tidak ditemukan di token",
					})
				return
			}

			if role != "admin" {
				context.AbortWithStatusJSON(http.StatusForbidden,
					gin.H{
						"error": "Akses ditolak, hanya admin yang diizinkan",
					})
				return
			}

			context.Set("role", role)
			context.Next()
		} else {
			context.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{
					"error": "Token tidak valid",
				})
		}
	}

}
