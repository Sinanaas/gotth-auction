package models

import (
	"time"

	"github.com/google/uuid"
)

type Auction struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Title string `gorm:"type:varchar(100)"`
	Description string `gorm:"type:text"`
	StartPrice float64 `gorm:"type:decimal(10,2)"`
	CurrentPrice float64 `gorm:"type:decimal(10,2)"`
	StartTime time.Time 
	EndTime time.Time 
	CreatedAt time.Time
	UpdatedAt time.Time
}