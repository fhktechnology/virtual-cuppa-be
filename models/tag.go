package models

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	OrganisationID uint           `gorm:"not null;index" json:"organisationId"`
	Name           string         `gorm:"type:varchar(100);not null" json:"name"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
