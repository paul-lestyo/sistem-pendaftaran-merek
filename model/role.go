package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID   uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name string    `gorm:"not null"`
}

type Roles struct {
	Roles []Role
}

func (role *Role) BeforeCreate(tx *gorm.DB) (err error) {
	role.ID = uuid.New()
	return
}
