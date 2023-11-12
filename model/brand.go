package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Brand struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	BusinessID  uuid.UUID
	Business    Business
	BrandName   string `gorm:"not null"`
	DescBrand   string `gorm:"not null"`
	BrandLogo   string `gorm:"not null"`
	Status      string `gorm:"not null;type:enum('OK', 'Perbaiki', 'Tolak', 'Menunggu')"`
	Note        string `gorm:"not null"`
	CreatedByID uuid.UUID
	CreatedBy   *User `gorm:"foreignKey:CreatedByID;constraint:OnDelete:CASCADE"`
	UpdatedByID uuid.UUID
	UpdatedBy   *User `gorm:"foreignKey:UpdatedByID;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time
}

type Brands struct {
	Brands []Brand
}

func (brand *Brand) BeforeCreate(tx *gorm.DB) (err error) {
	brand.ID = uuid.New()
	return
}
