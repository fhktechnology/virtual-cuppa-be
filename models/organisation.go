package models

import (
	"time"

	"gorm.io/gorm"
)

type Organisation struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Name       string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"name" binding:"required"`
	CompanyUrl string         `gorm:"type:varchar(500)" json:"companyUrl,omitempty"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
