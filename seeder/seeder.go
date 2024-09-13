package main

import (
	"log"
	"time"

	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/Sinanaas/gotth-auction/models"
	"github.com/Sinanaas/gotth-auction/utils"
	"github.com/google/uuid"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables")
	}

	initializers.ConnectDB(&config)
}

func main() {
	// Start a transaction
	tx := initializers.DB.Begin()

	// Ensure the transaction is rolled back if there's an error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Fatal("? Transaction failed and rolled back")
		}
	}()

	// Create a User
	hashedPassword, err := utils.HashPassword("habatusauda98")
	if err != nil {
		tx.Rollback()
		log.Fatal("? Could not hash password")
	}
	user := models.User{
		ID:       uuid.New(),
		Username: "sinanaas",
		Password: hashedPassword,
		Email:    "pudge842@gmail.com",
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		log.Fatal("? Could not create user")
	}

	// Create a Category
	category := models.Category{
		ID:          uuid.New(),
		Name:        "Electronics",
		Description: "Electronic items and gadgets",
	}
	if err := tx.Create(&category).Error; err != nil {
		tx.Rollback()
		log.Fatal("? Could not create category")
	}

	// Create an Auction
	auction := models.Auction{
		ID:           uuid.New(),
		Title:        "AUCTION BARU BOS",
		Description:  "Auction baru bos",
		StartPrice:   1000.00,
		CurrentPrice: 1000.00,
		UserID:       user.ID,
		Categories:   []*models.Category{&category},
		StartTime:    time.Now(),
		EndTime:      time.Now().Add(96 * time.Hour),
	}
	if err := tx.Create(&auction).Error; err != nil {
		tx.Rollback()
		log.Fatal("? Could not create auction")
	}

	// Create another Auction
	auction2 := models.Auction{
		ID:           uuid.New(),
		Title:        "Second Auction",
		Description:  "This is the second auction",
		StartPrice:   2000.00,
		CurrentPrice: 2000.00,
		UserID:       user.ID,
		Categories:   []*models.Category{&category},
		StartTime:    time.Now(),
		EndTime:      time.Now().Add(48 * time.Hour),
	}
	if err := tx.Create(&auction2).Error; err != nil {
		tx.Rollback()
		log.Fatal("? Could not create second auction")
	}

	// Create a Bid
	// bid := models.Bid{
	// 	ID:        uuid.New(),
	// 	AuctionID: auction.ID,
	// 	UserID:    user.ID,
	// 	BidAmount: 1300.00,
	// 	BidTime:   time.Now(),
	// }
	// if err := tx.Create(&bid).Error; err != nil {
	// 	tx.Rollback()
	// 	log.Fatal("? Could not create bid")
	// }

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Fatal("? Could not commit transaction")
	}

	log.Println("? Seed data created successfully")
}
