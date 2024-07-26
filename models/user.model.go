package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Username  string    `gorm:"type:varchar(100);unique_index" json:"username,omitempty"`
	Password  string    `gorm:"type:varchar(100)" json:"password,omitempty"`
	Email     string    `gorm:"type:varchar(100);unique_index" json:"email,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SignUpInput struct {
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Email     string    `json:"email,omitempty"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}