package controllers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/Sinanaas/gotth-auction/models"
	"github.com/Sinanaas/gotth-auction/toast"
	"github.com/Sinanaas/gotth-auction/utils"
	"github.com/gin-contrib/sessions"
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
		toast := toast.Danger("Invalid input: " + err.Error())
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid input: " + err.Error()})
		return
	}

	if payload.Email == "" || payload.Password == "" {
		toast := toast.Danger("Email and password are required")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Email and password are required"})
		return
	}

	if !validateEmail(payload.Email) {
		toast := toast.Danger("Invalid email address")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email address"})
		return
	}

	var user models.User
	result := ac.DB.Where("email = ?", payload.Email).First(&user)
	if result.Error != nil {
		toast := toast.Danger("Invalid email or password")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid email or password"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		toast := toast.Danger("Invalid email or password")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid email or password"})
		return
	}

	config, err := initializers.LoadConfig(".")
	if err != nil {
		toast := toast.Danger("Configuration error")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Configuration error"})
	}

	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID.String(), config.AccessTokenPrivateKey)
	if err != nil {
		toast := toast.Danger(err.Error())
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	refreshToken, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID.String(), config.RefreshTokenPrivateKey)
	if err != nil {
		toast := toast.Danger(err.Error())
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Set cookies with appropriate flags
	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	// session
	session := sessions.Default(ctx)
	session.Set("user_id", user.ID.String())
	if err := session.Save(); err != nil {
		toast := toast.Danger("Failed to save session")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to save session"})
		return
	}

	// Handle redirection
	ctx.Header("HX-Redirect", "/")
	ctx.Status(http.StatusOK)
}

func (ac AuthController) SignUp(ctx *gin.Context) {
	var payload models.SignUpInput

	if err := ctx.ShouldBind(&payload); err != nil {
		toast := toast.Danger("Invalid input: " + err.Error())
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid input: " + err.Error()})
		return
	}

	if payload.Email == "" || payload.Password == "" || payload.Username == "" {
		toast := toast.Danger("Email, username and password are required")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Email, username and password are required"})
		return
	}

	if !validateEmail(payload.Email) {
		toast := toast.Danger("Invalid email address")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email address"})
		return
	}

	var user models.User
	result := ac.DB.Where("email = ?", payload.Email).First(&user)
	if result.Error == nil {
		toast := toast.Danger("User with that email already exists")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
		return
	}

	result = ac.DB.Where("username = ?", payload.Username).First(&user)
	if result.Error == nil {
		toast := toast.Danger("User with that username already exists")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that username already exists"})
		return
	}

	if len(payload.Password) < 6 {
		toast := toast.Danger("Password must be at least 6 characters long")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Password must be at least 6 characters long"})
		return
	}

	if len(payload.Username) < 3 {
		toast := toast.Danger("Username must be at least 3 characters long")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Username must be at least 3 characters long"})
		return
	}

	if len(payload.Username) > 20 {
		toast := toast.Danger("Username must be at most 20 characters long")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Username must be at most 20 characters long"})
		return
	}

	if len(payload.Password) > 20 {
		toast := toast.Danger("Password must be at most 20 characters long")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Password must be at most 20 characters long"})
		return
	}

	if payload.Password != payload.ConfirmPassword {
		toast := toast.Danger("Passwords do not match")
		toast.SetHXTriggerHeader(ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		toast := toast.Danger("Error hashing password: " + err.Error())
		toast.SetHXTriggerHeader(ctx)
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

	result = ac.DB.Create(&newUser)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			toast := toast.Danger("User with that email already exists")
			toast.SetHXTriggerHeader(ctx)
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
		} else {
			toast := toast.Danger("Something went wrong: " + result.Error.Error())
			toast.SetHXTriggerHeader(ctx)
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something went wrong: " + result.Error.Error()})
		}
		return
	}

	// ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "User created successfully"})
	ctx.Header("HX-Redirect", "/login")
}

func (ac *AuthController) RefreshToken(ctx *gin.Context) {
	message := "could not refresh token"
	log.Println("refreshing token 1")
	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}
	log.Println("refreshing token 2")
	config, _ := initializers.LoadConfig(".")
	log.Println("refreshing token 3")
	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	log.Println("refreshing token 4")

	var user models.User
	result := ac.DB.First(&user, "id = ?", fmt.Sprint(sub))
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}
	log.Println("refreshing token 5")
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID.String(), config.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	log.Println("refreshing token 6")

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": access_token})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, false)

	session := sessions.Default(ctx)
	session.Clear()

	ctx.Header("HX-Redirect", "/login")
	ctx.Status(http.StatusOK)
}

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}
