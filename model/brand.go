package model

import "github.com/google/uuid"

type Brand struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	BusinessID  int
	Business    Business
	BrandName   string `gorm:"not null"`
	DescBrand   string `gorm:"not null"`
	CreatedByID int
	CreatedBy   User `gorm:"foreignKey:CreatedByID"`
	UpdatedByID int
	UpdatedBy   User `gorm:"foreignKey:UpdatedByID"`
}

type Brands struct {
	Brands []Brand
}

func (brand *Brand) name() {
	brand.ID = uuid.New()
	return
}
