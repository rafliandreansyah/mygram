package model

import (
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type Photo struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id" form:"id"`
	Title string `gorm:"not null" valid:"required~title is blank" json:"title" form:"title"`
	Caption string `json:"caption" form:"caption"`
	PhotoUrl string `gorm:"not null" valid:"required~photo is blank" json:"photo_url" form:"photo_url"`
	UserID uuid.UUID `gorm:"type:uuid;not null" valid:"required~user id required"json:"user_id" form:"user_id"`
	User User `gorm:"foreignKey:UserID" valid:"-"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`
}

func(p *Photo) BeforeCreate(tx *gorm.DB) error{
	var err error

	p.ID = uuid.NewV4()

	_, err = govalidator.ValidateStruct(p)
	if err != nil {
		return err
	}

	return nil
}