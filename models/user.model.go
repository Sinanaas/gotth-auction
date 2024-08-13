package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    Username  string    `gorm:"type:varchar(100);unique_index"`
    Password  string    `gorm:"type:varchar(100)"`
    Email     string    `gorm:"type:varchar(100);unique_index"`
    PhotoURL  string    `gorm:"type:varchar(255)"`
    // Bids relation
    Bids      []Bid     `gorm:"foreignKey:UserID"`    
    // Auctions relation
    Auctions  []Auction `gorm:"foreignKey:UserID"`   
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

type UserClient struct {
	Hub  *AuctionHub
	User *User
	Conn *websocket.Conn
	Send chan []byte
}
