package handlers

import (
	"fmt"

	"github.com/Sinanaas/gotth-auction/controllers"
	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/Sinanaas/gotth-auction/templates"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type GetProfileHandler struct{}

func NewGetProfileHandler() GetProfileHandler {
	return GetProfileHandler{}
}

func (gph GetProfileHandler) ServeHTTP(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var user_id string
	v := session.Get("user_id")
	if v != nil {
		user_id = v.(string)
	}
	fmt.Println("user_id: ", user_id)

	c := templates.Profile(controllers.NewBasicController(initializers.DB).GetUser(user_id))
	err := templates.Layout(c, user_id).Render(ctx.Request.Context(), ctx.Writer)

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
}
