package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/Sinanaas/gotth-auction/models"
	"github.com/Sinanaas/gotth-auction/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

func (ac AuthController) Login(ctx *gin.Context) {
	var payload models.SignInInput

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid input: " + err.Error()})
		return
	}

	var user models.User
	result := ac.DB.Where("email = ?", payload.Email).First(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid email or password"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid email or password"})
		return
	}

	config, err := initializers.LoadConfig(".")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Configuration error"})
		return
	}

	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID.String(), config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	refreshToken, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID.String(), config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Set cookies with appropriate flags
	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	// Handle redirection
	ctx.Header("HX-Redirect", "/")
	ctx.Status(http.StatusOK)
}

func (ac AuthController) SignUp(ctx *gin.Context) {
	var payload models.SignUpInput

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid input: " + err.Error()})
		return
	}

	if payload.Password != payload.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Error hashing password: " + err.Error()})
		return
	}

	now := time.Now()
	newUser := models.User{
		Email:     strings.ToLower(payload.Email),
		Password:  hashedPassword,
		Username:  payload.Username,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ac.DB.Create(&newUser)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
		} else {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something went wrong: " + result.Error.Error()})
		}
		return
	}

	// ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "User created successfully"})
	ctx.Header("HX-Redirect", "/login")
}

func (ac *AuthController) RefreshToken(ctx *gin.Context) {
	message := "could not refresh token"

	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config, _ := initializers.LoadConfig(".")

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "id = ?", fmt.Sprint(sub))
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID.String(), config.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": access_token})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
