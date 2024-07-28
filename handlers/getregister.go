package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-auction/templates"
	"github.com/gin-gonic/gin"
)

type GetRegisterHandler struct{}

func NewGetRegisterHandler() *GetRegisterHandler {
	return &GetRegisterHandler{}
}

func (h *GetRegisterHandler) ServeHTTP(ctx *gin.Context) {
	c := templates.Register("Register")
	err := templates.Layout(c, "My Website").Render(ctx.Request.Context(), ctx.Writer)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
