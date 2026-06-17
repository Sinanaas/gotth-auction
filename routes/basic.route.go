package routes

import (
	"github.com/Sinanaas/gotth-auction/controllers"
	"github.com/Sinanaas/gotth-auction/handlers"
	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/Sinanaas/gotth-auction/middleware"
	"github.com/gin-gonic/gin"
)

type BasicRouterController struct {
	basicController controllers.BasicController
}

func NewBasicRouterController(basicController controllers.BasicController) BasicRouterController {
	return BasicRouterController{basicController}
}

func (bc *BasicRouterController) BasicRoute(rg *gin.RouterGroup, config initializers.Config) {
	rg.GET("/", middleware.DeserializeUser(config), handlers.NewGetHomeHandler().ServeHTTP)
	rg.GET("/about", middleware.DeserializeUser(config), handlers.NewGetAboutHandler().ServeHTTP)
	rg.GET("/profile", middleware.DeserializeUser(config), handlers.NewGetProfileHandler().ServeHTTP)
	rg.POST("/update-profile", middleware.DeserializeUser(config), handlers.NewPostProfileHandler().ServeHTTP)
	rg.GET("/auction/:id", middleware.DeserializeUser(config), handlers.NewGetAuctionHandler().ServeHTTP)
	rg.GET("/history", middleware.DeserializeUser(config), handlers.NewGetHistoryHandler().ServeHTTP)
}
