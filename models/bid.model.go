package models

import (
	"time"

	"github.com/google/uuid"
)

type Bid struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	AuctionID uuid.UUID `gorm:"type:uuid;not null"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	BidAmount float64 `gorm:"type:decimal(10,2)"`
	BidTime time.Time
}