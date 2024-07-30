package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/Sinanaas/gotth-auction/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SaveFile saves the uploaded file, deletes the previous file if it exists, and returns the new file URL.
func SaveFile(ctx *gin.Context, file *multipart.FileHeader, userID string, db *gorm.DB) (string, error) {
	// Extract the file extension
	ext := filepath.Ext(file.Filename)

	// Create a new file name using the user ID
	newFileName := fmt.Sprintf("%s%s", userID, ext)

	// Define the file path
	filePath := filepath.Join("uploads", newFileName)

	// Fetch the user record from the database to get the old file URL
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return "", err
	}

	// Delete the old file if it exists and is different from the new file
	if user.PhotoURL != "" {
		oldFilePath := filepath.Join("uploads", user.PhotoURL)
		if oldFilePath != filePath {
			if err := os.Remove(oldFilePath); err != nil {
				return "", err
			}
		}
	}

	// Save the new file to the specified path
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		return "", err
	}

	// Return the new file URL
	return newFileName, nil
}
