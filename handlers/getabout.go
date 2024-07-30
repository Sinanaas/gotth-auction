package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-auction/templates"
	"github.com/gin-gonic/gin"
)

type GetAboutHandler struct{}

func NewGetAboutHandler() *GetAboutHandler {
	return &GetAboutHandler{}
}

func (gh *GetAboutHandler) ServeHTTP(ctx *gin.Context) {
	c := templates.About()
	err := templates.Layout(c, "").Render(ctx.Request.Context(), ctx.Writer)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
