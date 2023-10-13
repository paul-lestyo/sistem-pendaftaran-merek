package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Business struct {
	ID                uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID            uuid.UUID
	User              User
	BusinessName      string `gorm:"not null"`
	BusinessAddress   string `gorm:"not null"`
	OwnerName         string `gorm:"not null"`
	UMKCertificateUrl string `gorm:"not null"`
	SignatureUrl      string `gorm:"not null"`
}

type BusinessResponse struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	BusinessName      string
	BusinessAddress   string
	OwnerName         string
	UMKCertificateUrl string
	SignatureUrl      string
}

func (BusinessResponse) TableName() string {
	return "businesses"
}

type Businesses struct {
	Businesses []Business
}

func (business *Business) BeforeCreate(tx *gorm.DB) (err error) {
	business.ID = uuid.New()
	return
}
