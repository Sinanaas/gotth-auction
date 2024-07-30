package controllers

import (
	"github.com/Sinanaas/gotth-auction/models"
	"gorm.io/gorm"
)

type BasicController struct {
	DB *gorm.DB
}

func NewBasicController(DB *gorm.DB) BasicController {
	return BasicController{DB}
}

func (bc BasicController) GetUser(userId string) models.User {
	var user models.User
	result := bc.DB.Where("id = ?", userId).First(&user)
	if result.Error != nil {
		return models.User{}
	}
	return user
}
