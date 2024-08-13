package controllers

import (
	"log"
	"net/http"

	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/Sinanaas/gotth-auction/models"
	"github.com/Sinanaas/gotth-auction/toast"
	"github.com/Sinanaas/gotth-auction/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

func (bc BasicController) UpdateProfile(ctx *gin.Context) {
	// session
	session := sessions.Default(ctx)
	var dummy models.EditUserInput
	userID := session.Get("user_id")
	if userID == nil {
		toast := toast.Danger("Unauthorized")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := bc.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		toast := toast.Danger("User not found")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Bind form data to user model
	if err := ctx.ShouldBind(&dummy); err != nil {
		toast := toast.Danger("Invalid input")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if the username or email already exists
	var existingUser models.User
	if err := bc.DB.Where("username = ? OR email = ?", dummy.Username, dummy.Email).First(&existingUser).Error; err == nil && existingUser.ID != user.ID {
		toast := toast.Danger("Username or email already exists")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	if dummy.Username == "" {
		dummy.Username = user.Username
	}

	if dummy.Email == "" {
		dummy.Email = user.Email
	}

	// Handle file upload
	file, err := ctx.FormFile("profile_image")
	if err == nil {
		fileURL, err := utils.SaveFile(ctx, file, userID.(string), initializers.DB)
		if err != nil {
			toast := toast.Danger("Failed to upload file")
			toast.SetHXTriggerHeader(ctx)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
			return
		}
		dummy.PhotoURL = fileURL
	} else {
		dummy.PhotoURL = user.PhotoURL // Retain the current photo URL if no new photo is uploaded
	}

	// Update the user's profile
	if err := bc.DB.Model(&user).Where("id = ?", userID).Updates(map[string]interface{}{
		"username":  dummy.Username,
		"email":     dummy.Email,
		"photo_url": dummy.PhotoURL,
	}).Error; err != nil {
		toast := toast.Danger("Failed to update profile")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	toast := toast.Success("Profile updated successfully")
	toast.SetHXTriggerHeader(ctx)
}

func (bc BasicController) GetCategories() []models.Category {
	var categories []models.Category
	bc.DB.Find(&categories)
	return categories
}

func (bc BasicController) GetAuctions() []models.Auction {
	var auctions []models.Auction
	bc.DB.Model(&models.Auction{}).Preload("User").Preload("Category").Find(&auctions)
	return auctions
}

func (bc BasicController) GetAuction(auction_id string) models.Auction {
	var auction models.Auction
	result := bc.DB.Preload("User").Preload("Category").Preload("Bid").Where("id = ?", auction_id).First(&auction)
	if result.Error != nil {
		log.Printf("Error fetching auction: %v", result.Error)
	}
	return auction
}

func (bc BasicController) GetBidsForAuction(auctionID string) []models.Bid {
	var bids []models.Bid
	result := bc.DB.Preload("User").Where("auction_id = ?", auctionID).Order("bid_time desc").Find(&bids)
	if result.Error != nil {
		log.Printf("Error fetching bids: %v", result.Error)
	}
	return bids
}
