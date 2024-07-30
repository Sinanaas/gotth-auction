package handlers

import (
	"github.com/Sinanaas/gotth-auction/controllers"
	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/gin-gonic/gin"
)

type PostProfileHandler struct{}

func NewPostProfileHandler() *PostProfileHandler {
	return &PostProfileHandler{}
}

func (h *PostProfileHandler) ServeHTTP(ctx *gin.Context) {
	controllers.NewBasicController(initializers.DB).UpdateProfile(ctx)
}
