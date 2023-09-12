package model

import "github.com/google/uuid"

type Brand struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	BusinessID  uuid.UUID
	Business    Business
	BrandName   string `gorm:"not null"`
	DescBrand   string `gorm:"not null"`
	CreatedByID uuid.UUID
	CreatedBy   User `gorm:"foreignKey:CreatedByID"`
	UpdatedByID uuid.UUID
	UpdatedBy   User `gorm:"foreignKey:UpdatedByID"`
}

type Brands struct {
	Brands []Brand
}

func (brand *Brand) name() {
	brand.ID = uuid.New()
	return
}
