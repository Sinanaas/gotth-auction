package models

import "github.com/google/uuid"

type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Name        string    `gorm:"type:varchar(100);unique_index" json:"name,omitempty"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
}
