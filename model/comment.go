package model

import (
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id" form:"id"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" valid:"-"`
	Photo     Photo     `gorm:"foreignKey:PhotoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" valid:"-"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id" form:"user_id"`
	PhotoID   uuid.UUID `gorm:"type:uuid;not null" json:"photo_id" form:"photo_id"`
	Message   string    `json:"message" form:"message" valid:"required~message is blank"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`
}

func (comment *Comment) BeforeCreate(tx *gorm.DB) error{
	var err error

	comment.ID = uuid.NewV4()

	_, err = govalidator.ValidateStruct(comment)
	if err != nil {
		return err
	}

	return nil
}
