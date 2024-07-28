package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Username  string    `gorm:"type:varchar(100);unique_index"`
	Password  string    `gorm:"type:varchar(100)"`
	Email     string    `gorm:"type:varchar(100);unique_index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SignUpInput struct {
	Email           string
	Username		string
	Password        string 
	ConfirmPassword string 
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