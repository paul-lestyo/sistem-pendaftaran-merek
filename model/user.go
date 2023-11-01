package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name     string    `gorm:"type:varchar(255);not null"`
	Email    string    `gorm:"type:varchar(255);unique,not null"`
	Password string    `gorm:"type:varchar(255);not null"`
	ImageUrl string    `gorm:"type:varchar(255);not null"`
	RoleID   uuid.UUID
	Role     Role
	Business *Business `gorm:"constraint:OnDelete:CASCADE;"`
}

type Users struct {
	Users []User
}

type UsersBusinessCountBrands struct {
	ID           string
	Name         string
	Email        string
	ImageUrl     string
	BusinessName string
	CountBrands  int64
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	//	UUID version 4
	user.ID = uuid.New()
	return
}
