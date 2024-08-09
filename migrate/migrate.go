package main

import (
	"fmt"
	"log"

	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/Sinanaas/gotth-auction/models"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables")
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(&models.User{},  &models.Auction{}, &models.Category{}, &models.Bid{})
	fmt.Println("? Migration complete")
}
