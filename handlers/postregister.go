package handlers

import (
	"github.com/Sinanaas/gotth-auction/controllers"
	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/gin-gonic/gin"
)

type PostRegisterHandler struct{}

func NewPostRegisterHandler() *PostRegisterHandler {
	return &PostRegisterHandler{}
}

func (h *PostRegisterHandler) ServeHTTP(ctx *gin.Context) {
	controllers.NewAuthController(initializers.DB).SignUp(ctx)
}
