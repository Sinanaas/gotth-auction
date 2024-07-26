package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/Sinanaas/gotth-auction/controllers"
	"github.com/Sinanaas/gotth-auction/middleware"
)

type AuthRouterController struct {
	authController controllers.AuthController
}

func NewAuthRouterController(authController controllers.AuthController) AuthRouterController {
	return AuthRouterController{authController}
}

func (ac *AuthRouterController) AuthRoute(rg *gin.RouterGroup) {
	// rg.GET("/login", handlers.NewGetLoginHandler().ServeHTTP)
	// rg.POST("/login", handlers.NewPostLoginHandler().ServeHTTP)
	// rg.GET("/register", handlers.NewGetRegisterHandler().ServeHTTP)
	// rg.POST("/register", handlers.NewPostRegisterHandler().ServeHTTP)
	rg.GET("/logout", middleware.DeserializeUser(), ac.authController.LogoutUser)
}
