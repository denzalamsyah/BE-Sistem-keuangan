package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			if ctx.GetHeader("Content-Type") == "application/json" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "Unauthorized: Authorization header not found",
				})
				return
			} else {
				ctx.Redirect(http.StatusSeeOther, "/user/login")
			}
			return
		}

		// Mengambil token dari header Authorization
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(models.JwtKey), nil
		})
		if err != nil || !token.Valid {
			log.Printf("Error parsing or validating token: %v", err)
			if ctx.GetHeader("Content-Type") == "application/json" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "Unauthorized: Invalid session token",
				})
				return
			} else {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}

		ctx.Set("email", claims.Email)

		ctx.Next()
	})
}

