package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-auction/controllers"
	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/Sinanaas/gotth-auction/templates"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type GetAuctionHandler struct{}

func NewGetAuctionHandler() *GetAuctionHandler {
	return &GetAuctionHandler{}
}

func (ga *GetAuctionHandler) ServeHTTP(ctx *gin.Context) {
	bc := controllers.NewBasicController(initializers.DB)
	auction_id := ctx.Param("id")
	
	// get the auction and all the bids for that auction
	auction := bc.GetAuction(auction_id)
	bidders := bc.GetBidsForAuction(auction_id)
	
	c := templates.Auction(auction, bidders)
	session := sessions.Default(ctx)
	var user_id string
	v := session.Get("user_id")
	if v != nil {
		user_id = v.(string)
	}

	err := templates.Layout(c, user_id).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
