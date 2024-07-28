package handlers

import (
	"github.com/Sinanaas/gotth-auction/controllers"
	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/gin-gonic/gin"
)

type PostLoginHandler struct{}

func NewPostLoginHandler() *PostLoginHandler {
	return &PostLoginHandler{}
}

func (h *PostLoginHandler) ServeHTTP(ctx *gin.Context) {
	controllers.NewAuthController(initializers.DB).Login(ctx)
}
