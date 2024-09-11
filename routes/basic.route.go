package routes

import (
	"github.com/Sinanaas/gotth-auction/controllers"
	"github.com/Sinanaas/gotth-auction/handlers"
	"github.com/Sinanaas/gotth-auction/middleware"
	"github.com/gin-gonic/gin"
)

type BasicRouterController struct {
	basicController controllers.BasicController
}

func NewBasicRouterController(basicController controllers.BasicController) BasicRouterController {
	return BasicRouterController{basicController}
}

func (bc *BasicRouterController) BasicRoute(rg *gin.RouterGroup) {
	rg.GET("/", middleware.DeserializeUser(), handlers.NewGetHomeHandler().ServeHTTP)
	rg.GET("/about", middleware.DeserializeUser(), handlers.NewGetAboutHandler().ServeHTTP)
	rg.GET("/profile", middleware.DeserializeUser(), handlers.NewGetProfileHandler().ServeHTTP)
	rg.POST("/update-profile", middleware.DeserializeUser(), handlers.NewPostProfileHandler().ServeHTTP)
	rg.GET("/auction/:id", middleware.DeserializeUser(), handlers.NewGetAuctionHandler().ServeHTTP)
	rg.GET("/history", middleware.DeserializeUser(), handlers.NewGetHistoryHandler().ServeHTTP)
}
