package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name     string    `gorm:"not null"`
	Email    string    `gorm:"unique,not null"`
	Password string    `gorm:"not null"`
	ImageUrl string    `gorm:"not null"`
	RoleID   int
	Role     Role
}

type Users struct {
	Users []User
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	//	UUID version 4
	user.ID = uuid.New()
	return
}
