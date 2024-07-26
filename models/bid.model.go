package models

import (
	"time"

	"github.com/google/uuid"
)

type Bid struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	AuctionID uuid.UUID `gorm:"type:uuid;not null" json:"auction_id,omitempty"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id,omitempty"`
	BidAmount float64 `gorm:"type:decimal(10,2)" json:"bid_amount,omitempty"`
	BidTime time.Time `json:"bid_time,omitempty"`
}