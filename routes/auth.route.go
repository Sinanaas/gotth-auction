package routes

import (
	"github.com/Sinanaas/gotth-auction/controllers"
	"github.com/Sinanaas/gotth-auction/handlers"
	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/Sinanaas/gotth-auction/middleware"
	"github.com/gin-gonic/gin"
)

type AuthRouterController struct {
	authController controllers.AuthController
}

func NewAuthRouterController(authController controllers.AuthController) AuthRouterController {
	return AuthRouterController{authController}
}

func (ac *AuthRouterController) AuthRoute(rg *gin.RouterGroup, config initializers.Config) {
	rg.GET("/login", handlers.NewGetLoginHandler().ServeHTTP)
	rg.POST("/login", handlers.NewPostLoginHandler().ServeHTTP)
	rg.GET("/register", handlers.NewGetRegisterHandler().ServeHTTP)
	rg.POST("/register", handlers.NewPostRegisterHandler().ServeHTTP)
	rg.GET("/logout", middleware.DeserializeUser(config), ac.authController.LogoutUser)
	rg.GET("/refresh", middleware.DeserializeUser(config), ac.authController.RefreshToken)
}
