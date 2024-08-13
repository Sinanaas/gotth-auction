package routes

import (
	"github.com/Sinanaas/gotth-auction/controllers"
	"github.com/Sinanaas/gotth-auction/middleware"
	"github.com/gin-gonic/gin"
)

type WebsocketRouterController struct {
	websocketcController controllers.WebsocketController
}

func NewWebsocketRouterController(WebsocketController controllers.WebsocketController) WebsocketRouterController {
	return WebsocketRouterController{WebsocketController}
}

func (wc *WebsocketRouterController) WebsocketRoute(rg *gin.RouterGroup) {
	hub := controllers.NewAuctionHub()
	go controllers.Run(hub)
	rg.GET("/ws", middleware.DeserializeUser(), func(ctx *gin.Context) {
		wc.websocketcController.ServeWS(hub, ctx)
	})
}
