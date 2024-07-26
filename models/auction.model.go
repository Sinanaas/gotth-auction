package models

import (
	"time"

	"github.com/google/uuid"
)

type Auction struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Title string `gorm:"type:varchar(100)" json:"title,omitempty"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	StartPrice float64 `gorm:"type:decimal(10,2)" json:"start_price,omitempty"`
	CurrentPrice float64 `gorm:"type:decimal(10,2)" json:"current_price,omitempty"`
	StartTime time.Time `json:"start_time,omitempty"`
	EndTime time.Time `json:"end_time,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}