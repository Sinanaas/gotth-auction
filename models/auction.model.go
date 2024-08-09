package models

import (
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
	UserID       uuid.UUID   `gorm:"type:uuid;not null;"`
	User         User        `gorm:"foreignKey:UserID"`
	Categories   []*Category `gorm:"many2many:auction_categories;"`
	PhotoURL     string      `gorm:"type:varchar(255)"`
	StartTime    time.Time
	EndTime      time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// type AuctionHub struct {
// 	clients map[*models.User]bool
// 	broadcast chan []byte
// 	register chan *models.User
// 	unregister chan *models.User
// }
