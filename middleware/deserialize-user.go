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

func DeserializeUser(config initializers.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		access_token, err := ctx.Cookie("access_token")
		if err != nil || access_token == "" {
			ctx.Redirect(http.StatusSeeOther, "/login")
			return
		}

		sub, err := utils.ValidateToken(access_token, config.AccessTokenPublicKey)
		if err != nil {
			at := controllers.NewAuthController(initializers.DB)
			at.RefreshToken(ctx)
			if ctx.IsAborted() {
				return
			}

			access_token, _ = ctx.Cookie("access_token")
			sub, err = utils.ValidateToken(access_token, config.AccessTokenPublicKey)
			if err != nil {
				ctx.Redirect(http.StatusSeeOther, "/login")
				return
			}
		}

		var user models.User
		if result := initializers.DB.First(&user, "id = ?", fmt.Sprint(sub)); result.Error != nil {
			ctx.Redirect(http.StatusSeeOther, "/login")
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
