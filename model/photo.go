package model

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Photo struct {
	ID uuid.UUID `gorm:"primaryKey" json:"id" form:"id"`
	Title string `gorm:"not null" valid:"required~title is blank" json:"title" form:"title"`
	Caption string `json:"caption" form:"caption"`
	PhotoUrl string `gorm:"not null" valid:"required~photo is blank" json:"photo_url" form:"photo_url"`
	UserID uuid.UUID `gorm:"not null" valid:"required~user id required"json:"user_id" form:"user_id"`
	User User `gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`
}