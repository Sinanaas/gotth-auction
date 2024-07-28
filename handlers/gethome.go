package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-auction/templates"
	"github.com/gin-gonic/gin"
)

type GetHomeHandler struct{}

func NewGetHomeHandler() *GetHomeHandler {
	return &GetHomeHandler{}
}

func (gh *GetHomeHandler) ServeHTTP(ctx *gin.Context) {
	c := templates.Home()
	err := templates.Layout(c, "My Website").Render(ctx.Request.Context(), ctx.Writer)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
