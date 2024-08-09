package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Username  string    `gorm:"type:varchar(100);unique_index"`
	Password  string    `gorm:"type:varchar(100)"`
	Email     string    `gorm:"type:varchar(100);unique_index"`
	PhotoURL  string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SignUpInput struct {
	Email           string
	Username        string
	Password        string
	ConfirmPassword string
}

type EditUserInput struct {
	Username string
	Email    string
	PhotoURL string
}

type SignInInput struct {
	Email    string
	Password string
}

type UserResponse struct {
	ID        uuid.UUID
	Email     string
	Username  string
	Provider  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// type UserClient struct {
// 	hub  *AuctionHub
// 	user *User
// 	conn *websocket.Conn
// 	send chan []byte
// }

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// func serveWS(hub *AuctionHub, ctx *gin.Context) {
// 	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// }
