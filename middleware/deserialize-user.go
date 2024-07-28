package middleware

import (
	"fmt"
	"net/http"

	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/Sinanaas/gotth-auction/models"
	"github.com/Sinanaas/gotth-auction/utils"
	"github.com/gin-gonic/gin"
)

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var access_token string
		cookie, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
			return
		}

		access_token = cookie
		if access_token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
			return
		}

		config, _ := initializers.LoadConfig(".")
		sub, err := utils.ValidateToken(access_token, config.AccessTokenPublicKey)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "Aerror", "message": err.Error()})
			return
		}

		var user models.User
		result := initializers.DB.First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
