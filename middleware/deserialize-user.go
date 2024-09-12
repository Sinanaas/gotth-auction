package middleware

import (
	"fmt"
	"net/http"

	"github.com/Sinanaas/gotth-auction/controllers"
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
			ctx.Redirect(http.StatusSeeOther, "/login")
			return
		}

		access_token = cookie
		if access_token == "" {
			ctx.Redirect(http.StatusSeeOther, "/login")
			return
		}

		config, _ := initializers.LoadConfig(".")
		sub, err := utils.ValidateToken(access_token, config.AccessTokenPublicKey)

		if err != nil {
			ctx.Redirect(http.StatusSeeOther, "/login")
			return
		}

		var user models.User
		result := initializers.DB.First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			ctx.Redirect(http.StatusSeeOther, "/login")
			return
		}
		
		var at = controllers.NewAuthController(initializers.DB)
		at.RefreshToken(ctx)

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
