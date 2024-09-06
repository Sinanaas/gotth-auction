package routes

import (
	"github.com/Sinanaas/gotth-auction/controllers"
	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/Sinanaas/gotth-auction/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WebsocketRouterController struct {
	websocketController controllers.WebsocketController
}

func NewWebsocketRouterController(WebsocketController controllers.WebsocketController) WebsocketRouterController {
	return WebsocketRouterController{WebsocketController}
}

func (wc *WebsocketRouterController) WebsocketRoute(rg *gin.RouterGroup) {
	hub := controllers.NewAuctionHub()
	
	rg.GET("/ws/:id", middleware.DeserializeUser(), func(ctx *gin.Context) {
		auction_id := ctx.Param("id")
		
		auction := controllers.NewBasicController(initializers.DB).GetAuction(auction_id)
		
		hub.Auction = &auction
		
		parsed_auction_id, err := uuid.Parse(auction_id)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid auction ID"})
			return
		}
		hub.Auction.ID = parsed_auction_id
		
		go controllers.Run(hub)
		wc.websocketController.ServeWS(hub, ctx)
	})
}
