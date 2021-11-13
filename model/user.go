package model

import (
	"MyGram/helper"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id" form:"id"`
	UserName  string    `gorm:"not null;unique" valid:"required~Username is blank" json:"username" form:"username"`
	Email     string    `gorm:"not null;unique" valid:"required~Email is blank,email~Invalid email format" json:"email" form:"email"`
	Password  string    `gorm:"not null" valid:"required~Password is blank,minstringlength(8)~Password must have 8 character or more" json:"password" form:"password"`
	Age       int       `gorm:"not null" valid:"required~Age is blank" json:"age" form:"age"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	var err error
	u.ID = uuid.NewV4()

	_, err = govalidator.ValidateStruct(u)
	if err != nil {
		return err
	}

	hashPass, err := helper.HashPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = hashPass
	return nil
}
