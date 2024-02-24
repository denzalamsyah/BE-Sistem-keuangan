package middleware

import (
	"net/http"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		// Mengambil nilai session token dari cookie
	sessionToken, err := ctx.Cookie("session_token")
	if err != nil {
		if ctx.GetHeader("Content-Type") == "application/json" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized: Session token not found",
			})
			return
		} else {
			ctx.Redirect(http.StatusSeeOther, "/user/login")
		}
		return
	}

	// fmt.Println("Session Token:", sessionToken)
	
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(sessionToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(models.JwtKey), nil
	})
	// if err != nil || !token.Valid {
	// 	log.Println("Error parsing or validating token:", err)
	// 	if ctx.GetHeader("Content-Type") == "application/json" {
	// 		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
	// 			"message": "Unauthorized: Invalid session token",
	// 		})
	// 		return
	// 	} else {
	// 		ctx.AbortWithStatus(http.StatusUnauthorized)
	// 		return
	// 	}
	// }

	if err != nil || !token.Valid {
		ctx.AbortWithStatus(400)
		return
	}

	ctx.Set("email", claims.Email)

	ctx.Next()
	})
}
