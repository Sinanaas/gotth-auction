package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
    gorm.Model
    ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    Name        string     `gorm:"type:varchar(100);unique_index"`
    Description string     `gorm:"type:text"`
    // Auctions relation
    Auctions    []*Auction `gorm:"many2many:auction_categories;"`
}
