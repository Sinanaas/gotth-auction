package routes

import (
	"time"

	"github.com/Sinanaas/gotth-auction/controllers"
	"github.com/Sinanaas/gotth-auction/handlers"
	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/Sinanaas/gotth-auction/middleware"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
)

type AuthRouterController struct {
	authController controllers.AuthController
}

func NewAuthRouterController(authController controllers.AuthController) AuthRouterController {
	return AuthRouterController{authController}
}

func (ac *AuthRouterController) AuthRoute(rg *gin.RouterGroup, config initializers.Config) {
	loginRate := limiter.Rate{Period: 1 * time.Minute, Limit: 10}
	rg.GET("/login", handlers.NewGetLoginHandler().ServeHTTP)
	rg.POST("/login", middleware.RateLimiter(loginRate), handlers.NewPostLoginHandler().ServeHTTP)
	rg.GET("/register", handlers.NewGetRegisterHandler().ServeHTTP)
	rg.POST("/register", middleware.RateLimiter(loginRate), handlers.NewPostRegisterHandler().ServeHTTP)
	rg.GET("/logout", middleware.DeserializeUser(config), ac.authController.LogoutUser)
	rg.GET("/refresh", middleware.DeserializeUser(config), ac.authController.RefreshToken)
}
