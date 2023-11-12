package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Announcement struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	Title     string    `gorm:"not null"`
	Desc      string    `gorm:"type:text; not null"`
	Tag       string    `gorm:"not null"`
	ImageUrl  string
	CreatedBy uuid.UUID
	User      User `gorm:"foreignKey:CreatedBy"`
	CreatedAt *time.Time
}

type Announcements struct {
	Announcements []Announcement
}

func (announcement *Announcement) BeforeCreate(tx *gorm.DB) (err error) {
	announcement.ID = uuid.New()
	return
}
