package controllers

import "gorm.io/gorm"

type BasicController struct {
	DB *gorm.DB
}

func NewBasicController(DB *gorm.DB) BasicController {
	return BasicController{DB}
}
