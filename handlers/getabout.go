package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-auction/templates"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type GetAboutHandler struct{}

func NewGetAboutHandler() *GetAboutHandler {
	return &GetAboutHandler{}
}

func (gh *GetAboutHandler) ServeHTTP(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var user_id string
	v := session.Get("user_id")
	if v != nil {
		user_id = v.(string)
	}

	c := templates.About()
	err := templates.Layout(c, user_id).Render(ctx.Request.Context(), ctx.Writer)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
