package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-auction/controllers"
	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/Sinanaas/gotth-auction/templates"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type GetHistoryHandler struct{}

func NewGetHistoryHandler() *GetHistoryHandler {
	return &GetHistoryHandler{}
}

func (gh *GetHistoryHandler) ServeHTTP(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var user_id string
	v := session.Get("user_id")
	if v != nil {
		user_id = v.(string)
	}
	if user_id == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c := templates.History(controllers.NewBasicController(initializers.DB).GetHistory(user_id))
	err := templates.Layout(c, user_id).Render(ctx.Request.Context(), ctx.Writer)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
