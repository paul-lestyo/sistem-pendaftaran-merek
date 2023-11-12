package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Log struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID
	User      User
	CreatedAt *time.Time
}

type Logs struct {
	Logs []Log
}

func (log *Log) BeforeCreate(tx *gorm.DB) (err error) {
	log.ID = uuid.New()
	return
}
