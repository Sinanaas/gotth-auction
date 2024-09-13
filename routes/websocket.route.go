package routes

import (
	"github.com/Sinanaas/gotth-auction/controllers"
	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/Sinanaas/gotth-auction/middleware"
	"github.com/Sinanaas/gotth-auction/models"
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
	var autionHub []*models.AuctionHub
	
	rg.GET("/ws/:id", middleware.DeserializeUser(), func(ctx *gin.Context) {
		auction_id := ctx.Param("id")
		auction := controllers.NewBasicController(initializers.DB).GetAuction(auction_id)
		for _, h := range autionHub {
			if h.Auction.ID == auction.ID {
				go controllers.Run(h)
				wc.websocketController.ServeWS(h, ctx)
				return
			}
		}

		hub := controllers.NewAuctionHub()
		hub.Auction = &auction
		
		parsed_auction_id, err := uuid.Parse(auction_id)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid auction ID"})
			return
		}
		hub.Auction.ID = parsed_auction_id
		autionHub = append(autionHub, hub)
		go controllers.Run(hub)
		wc.websocketController.ServeWS(hub, ctx)
	})
}
