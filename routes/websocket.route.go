package routes

import (
	"time"

	"github.com/Sinanaas/gotth-auction/controllers"
	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/Sinanaas/gotth-auction/middleware"
	"github.com/Sinanaas/gotth-auction/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ulule/limiter/v3"
)

type WebsocketRouterController struct {
	websocketController controllers.WebsocketController
}

func NewWebsocketRouterController(WebsocketController controllers.WebsocketController) WebsocketRouterController {
	return WebsocketRouterController{WebsocketController}
}

func (wc *WebsocketRouterController) WebsocketRoute(rg *gin.RouterGroup, config initializers.Config) {
	var autionHub []*models.AuctionHub
	wsRate := limiter.Rate{Period: 1 * time.Minute, Limit: 10}

	rg.GET("/ws/:id", middleware.RateLimiter(wsRate), middleware.DeserializeUser(config), func(ctx *gin.Context) {
		auction_id := ctx.Param("id")
		auction := controllers.NewBasicController(initializers.DB).GetAuction(auction_id)
		for _, h := range autionHub {
			if h.Auction.ID == auction.ID {
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
