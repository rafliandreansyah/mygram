package model

import (
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type SocialMedia struct {
	ID             uuid.UUID `json:"id" form:"id" gorm:"type:uuid;primaryKey"`
	Name           string    `json:"name" form:"name" valid:"required~name is blank"`
	SocialMediaUrl string    `json:"social_media_url" form:"social_media_url" valid:"required~social media url is blank"`
	User           User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" valid:"-"`
	UserID         uuid.UUID `json:"user_id" form:"user_id"`
	CreatedAt      time.Time `json:"created_at" form:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" form:"updated_at"`
}

func (p *SocialMedia) BeforeCreate(tx *gorm.DB) error {
	var err error

	p.ID = uuid.NewV4()

	_, err = govalidator.ValidateStruct(p)
	if err != nil {
		return err
	}

	return nil
}
