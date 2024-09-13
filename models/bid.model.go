package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bid struct {
    gorm.Model
    ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    // Auction relation
    AuctionID uuid.UUID `gorm:"type:uuid;not null"`
    Auction   Auction   `gorm:"foreignKey:AuctionID"`  
    // User relation
    UserID    uuid.UUID `gorm:"type:uuid;not null;"`
    User      User      `gorm:"foreignKey:UserID"`      
    BidAmount float64   `gorm:"type:decimal(10,2)"`
    BidTime   time.Time
}

type WSBid struct {
	Price string `json:"price"`
	Headers interface{} `json:"headers"`
}
