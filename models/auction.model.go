package models

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Auction struct {
    gorm.Model
    ID           uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    Title        string      `gorm:"type:varchar(100)"`
    Description  string      `gorm:"type:text"`
    StartPrice   float64     `gorm:"type:decimal(10,2)"`
    CurrentPrice float64     `gorm:"type:decimal(10,2)"`
    // User relation
    UserID       uuid.UUID   `gorm:"type:uuid;not null;"`
    User         User        `gorm:"foreignKey:UserID"`         
    // Bids relation 
    Bids         []Bid       `gorm:"foreignKey:AuctionID"`      
    // Categories relation
    Categories   []*Category `gorm:"many2many:auction_categories;"`
    Winner       uuid.UUID   `gorm:"type:uuid;"`
    PhotoURL     string      `gorm:"type:varchar(255)"`
    StartTime    time.Time
    EndTime      time.Time
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type AuctionHub struct {
	sync.RWMutex
    
	Clients    map[*UserClient]bool
	Messages   []*Bid
	Broadcast  chan *Bid
	Auction    *Auction
	Register   chan *UserClient
	Unregister chan *UserClient
}
